package etcd

import (
	"context"
	"fmt"
	"testing"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func Test_Client(t *testing.T) {
	client, err := New(
		WithHost("etcd-server.liushuojia.com:2379"),
		WithDialTimeout(3*time.Second),
	)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer client.Close() // 退出前关闭连接

	// 3. 核心操作：增（Put）、查（Get）、改（Put 覆盖）、删（Delete）
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// ---------------------- 增/改：Put 键值对 ----------------------
	putResp, err := client.Put(ctx, "/config/user/max_count", "100")
	if err != nil {
		fmt.Printf("Put 失败：%v\n", err)
		return
	}
	fmt.Printf("Put 成功，Revision：%d\n", putResp.Header.Revision)

	// ---------------------- 查：Get 键值对 ----------------------
	// 单个键查询
	getResp, err := client.Get(ctx, "/config/user/max_count")
	if err != nil {
		fmt.Printf("Get 失败：%v\n", err)
		return
	}
	// 遍历结果（Get 可返回多个键，如前缀查询）
	for _, kv := range getResp.Kvs {
		fmt.Printf("键：%s，值：%s，版本：%d\n", kv.Key, kv.Value, kv.ModRevision)
	}

	// 前缀查询（常用：批量获取配置）
	prefixResp, err := client.Get(ctx, "/config/", clientv3.WithPrefix())
	if err != nil {
		fmt.Printf("前缀查询失败：%v\n", err)
		return
	}
	fmt.Println("\n前缀 /config/ 的所有键值：")
	for _, kv := range prefixResp.Kvs {
		fmt.Printf("  %s = %s\n", kv.Key, kv.Value)
	}

	// ---------------------- 删：Delete 键值对 ----------------------
	delResp, err := client.Delete(ctx, "/config/user/max_count")
	if err != nil {
		fmt.Printf("Delete 失败：%v\n", err)
		return
	}
	fmt.Printf("\nDelete 成功，删除数量：%d\n", delResp.Deleted)
}
func Test_Client_Watcher(t *testing.T) {
	client, err := New(
		WithHost("etcd-server.liushuojia.com:2379"),
		WithDialTimeout(3*time.Second),
	)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer client.Close() // 退出前关闭连接

	// 2. 创建 Watcher（监听 /config/ 前缀的所有键）
	watcher := client.Watch(context.Background(), "/config/", clientv3.WithPrefix())

	// 3. 持续监听（阻塞式，直到上下文取消）
	fmt.Println("开始监听 /config/ 前缀的键变更...")
	for watchResp := range watcher {
		// 遍历本次监听的所有事件
		for _, event := range watchResp.Events {
			fmt.Printf(
				"事件类型：%s，键：%s，值：%s，Revision：%d\n",
				event.Type,     // PUT/DELETE
				event.Kv.Key,   // 触发事件的键
				event.Kv.Value, // 触发事件的值
				event.Kv.ModRevision,
			)
		}
	}
}
func Test_Client_Lease(t *testing.T) {
	client, err := New(
		WithHost("etcd-server.liushuojia.com:2379"),
		WithDialTimeout(3*time.Second),
	)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer client.Close() // 退出前关闭连接

	// 2. 创建租约：TTL=10秒（10秒内不续期则租约过期）
	leaseCtx, leaseCancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer leaseCancel()
	leaseResp, err := client.Grant(leaseCtx, 10)
	if err != nil {
		panic(fmt.Sprintf("创建租约失败：%v", err))
	}
	leaseID := leaseResp.ID
	fmt.Printf("创建租约成功，LeaseID：%d，TTL：%d秒\n", leaseID, leaseResp.TTL)

	// 3. 绑定租约到键：创建临时节点（服务注册示例）
	putCtx, putCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer putCancel()
	_, err = client.Put(putCtx, "/config/192.168.1.100:8080", "user-service-v1", clientv3.WithLease(leaseID))
	if err != nil {
		panic(fmt.Sprintf("绑定租约失败：%v", err))
	}
	fmt.Println("创建临时节点成功：/config/192.168.1.100:8080")

	// 4. 自动续期（核心：防止租约过期，服务注册必须续期）
	keepAliveChan, err := client.KeepAlive(context.Background(), leaseID)
	if err != nil {
		panic(fmt.Sprintf("续期失败：%v", err))
	}
	// 监听续期结果（可选）
	go func() {
		for ka := range keepAliveChan {
			fmt.Printf("租约续期成功，剩余TTL：%d秒\n", ka.TTL)
		}
		fmt.Println("租约续期通道关闭（租约过期/手动取消）")
	}()

	go func() {
		time.Sleep(18 * time.Second)
		// 6. 手动撤销租约（可选，服务正常下线时主动注销）
		revokeCtx, revokeCancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer revokeCancel()
		_, err = client.Revoke(revokeCtx, leaseID)
		if err != nil {
			panic(fmt.Sprintf("撤销租约失败：%v", err))
		}
		fmt.Println("手动撤销租约成功")
	}()

	// 5. 模拟服务运行（15秒后退出，租约不再续期）
	fmt.Println("服务运行中，15秒后退出...")
	time.Sleep(30 * time.Second)

}
