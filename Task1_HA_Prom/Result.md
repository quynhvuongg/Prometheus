

**Task 1**

- Dựng Prom cluster 2 node và Alertmanager cluster 2 node bằng docker compose.
- kịch bản test: stop bất kì 1 node Prom và 1 node Alertmanager , Prom vẫn đảm bảo chạy bình thường , đồng thời có cảnh báo service down

**Alertmanager Cluster Status**


![ ](https://scontent.xx.fbcdn.net/v/t1.15752-0/p280x280/94068827_232707258016163_3663900930790653952_n.png?_nc_cat=106&_nc_sid=b96e70&_nc_ohc=RXRVFl2QtAcAX-H_bmm&_nc_ad=z-m&_nc_cid=0&_nc_zor=9&_nc_ht=scontent.xx&oh=e7dfa13ff1f6bbb596181fd0dd70bc4b&oe=5EC8D901)

![ ](https://scontent.xx.fbcdn.net/v/t1.15752-0/p280x280/94212798_678143253024719_8563699119279833088_n.png?_nc_cat=102&_nc_sid=b96e70&_nc_ohc=yXJdbcyE6U8AX_9OsJ6&_nc_ad=z-m&_nc_cid=0&_nc_zor=9&_nc_ht=scontent.xx&oh=0a65d249b86d22f51b0cac8387f50e2c&oe=5ECB3E5D)


**Kịch bản 1: Prometheus2 down**


![ ](https://scontent.xx.fbcdn.net/v/t1.15752-0/p280x280/93929269_253780952474624_174008826175946752_n.png?_nc_cat=110&_nc_sid=b96e70&_nc_ohc=LWVDaAiegPgAX_5Swym&_nc_ad=z-m&_nc_cid=0&_nc_zor=9&_nc_ht=scontent.xx&oh=e7b291185a06496d658c73504bf3b661&oe=5ECA63D0)

![ ](https://scontent.xx.fbcdn.net/v/t1.15752-0/p280x280/93964409_564954024143809_5459334914220490752_n.png?_nc_cat=105&_nc_sid=b96e70&_nc_ohc=P4c778KKZJIAX9_GaNc&_nc_ad=z-m&_nc_cid=0&_nc_zor=9&_nc_ht=scontent.xx&oh=61b5611d57d9191170cde476b7edd8ba&oe=5EC95986)

![ ](https://scontent.xx.fbcdn.net/v/t1.15752-0/p280x280/94073908_875784719590940_4745467785772007424_n.png?_nc_cat=108&_nc_sid=b96e70&_nc_ohc=J-INLtFGEsoAX-A3-Za&_nc_ad=z-m&_nc_cid=0&_nc_zor=9&_nc_ht=scontent.xx&oh=eef977f231c5bc84a1f3a7598d6da7ca&oe=5EC95FDD)

![ ](https://scontent.xx.fbcdn.net/v/t1.15752-0/p280x280/94568072_541167896582084_2186552702043947008_n.png?_nc_cat=110&_nc_sid=b96e70&_nc_ohc=l9D4KjChr7QAX9WZza_&_nc_ad=z-m&_nc_cid=0&_nc_zor=9&_nc_ht=scontent.xx&oh=7e7bbdf95d345e7a361fdaed950409dc&oe=5EC836E0)


**Kịch bản 2: Alertmanager1 down**


![ ](https://scontent.xx.fbcdn.net/v/t1.15752-0/p280x280/93961215_279797269692674_8532708924489465856_n.png?_nc_cat=108&_nc_sid=b96e70&_nc_ohc=kVRhHrfgKSsAX8Hf0rR&_nc_ad=z-m&_nc_cid=0&_nc_zor=9&_nc_ht=scontent.xx&oh=d94f90f0f2bacb352e044470352c1a2f&oe=5EC9A624)

![ ](https://scontent.xx.fbcdn.net/v/t1.15752-0/p280x280/94488833_231502138264752_7049108949936635904_n.png?_nc_cat=105&_nc_sid=b96e70&_nc_ohc=vlqveD0YLD4AX_ikWaw&_nc_ad=z-m&_nc_cid=0&_nc_zor=9&_nc_ht=scontent.xx&oh=40a1984946d9c914cd983629cbaec8da&oe=5EC816DC)

![ ](https://scontent.xx.fbcdn.net/v/t1.15752-0/p280x280/94006598_2957057271026321_8780454694095945728_n.png?_nc_cat=102&_nc_sid=b96e70&_nc_ohc=JJIin4OPK-AAX9KPgb8&_nc_ad=z-m&_nc_cid=0&_nc_zor=9&_nc_ht=scontent.xx&oh=81cb1321930700017e848a99666f6f3a&oe=5EC8B143)

![ ](https://scontent.xx.fbcdn.net/v/t1.15752-0/p280x280/94385129_266258924775231_7910562729612541952_n.png?_nc_cat=107&_nc_sid=b96e70&_nc_ohc=PV8z4NkajC4AX9VaEYa&_nc_ad=z-m&_nc_cid=0&_nc_zor=9&_nc_ht=scontent.xx&oh=fa809b496fcce2aaf8417f64a1ebef78&oe=5ECB0F94)





