**Task 2**

Federation: Dựng 2 node Prom X và Y. Prom X cấu hình targer giám sát host dùng node_exporter. Prom Y sử dụng federation để collect chỉ các metric về CPU từ Prom A.

* Targets trên 2 Prometheus:

![ ](https://scontent.xx.fbcdn.net/v/t1.15752-0/p280x280/94228504_276306013379884_1906658741449654272_n.png?_nc_cat=101&_nc_sid=b96e70&_nc_ohc=vrIbMo8uD1sAX96dsD3&_nc_ad=z-m&_nc_cid=0&_nc_zor=9&_nc_ht=scontent.xx&oh=d0275fecda9f6c97d9e859bfba6dbfa3&oe=5EC8F802)

![ ](https://scontent.xx.fbcdn.net/v/t1.15752-0/p280x280/94123576_570187246955627_3791953679913844736_n.png?_nc_cat=108&_nc_sid=b96e70&_nc_ohc=AY6f7LILBTIAX_4v4dy&_nc_ad=z-m&_nc_cid=0&_nc_zor=9&_nc_ht=scontent.xx&oh=06a2c47e8fc1cbbae3a092d15d9d0c32&oe=5ECA9692)

* node_cpu_seconds_total metrics

![ ](https://scontent.xx.fbcdn.net/v/t1.15752-0/p280x280/94132044_596848617587202_4088073502966415360_n.png?_nc_cat=104&_nc_sid=b96e70&_nc_ohc=DGkYwzKsKFQAX-8um5s&_nc_ad=z-m&_nc_cid=0&_nc_zor=9&_nc_ht=scontent.xx&oh=35ba6cac8baae92a0493dcda45a4d258&oe=5ECAE247)

![ ](https://scontent.xx.fbcdn.net/v/t1.15752-0/p280x280/94203371_248987966475951_6409427838109220864_n.png?_nc_cat=106&_nc_sid=b96e70&_nc_ohc=MgjjxPPej8oAX-Qr7Ad&_nc_ad=z-m&_nc_cid=0&_nc_zor=9&_nc_ht=scontent.xx&oh=d2fc155fe160089c963738c627e5bac8&oe=5EC8D4C1)
