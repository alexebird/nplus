nplus
=====

```
$nplus tasks
JOB                                          TASK              ASTATUS   TSTATE   TFAILED  JTYPE    ADDR                         ALLOCID   NODEID
ci-build-image/dispatch-1512599428-7e364f7a  ci-build-image    complete  dead     false    batch                                 760573ef  2fc75e14
ci-build-image/dispatch-1512598651-3b83b6b1  ci-build-image    complete  dead     false    batch                                 26012921  2fc75e14
ci-build-image/dispatch-1512586589-d5a24b14  ci-build-image    complete  dead     false    batch                                 163a703a  2fc75e14
ci-build-image/dispatch-1512586465-4221e47b  ci-build-image    complete  dead     false    batch                                 1cf03324  2fc75e14
ci-build-image/dispatch-1512586367-91459c4d  ci-build-image    complete  dead     false    batch                                 0c422b59  2fc75e14
infragc                                      wkr-ci            running   running  false    system   10.40.5.247:29119(http)      1118e083  2fc75e14
consul-exporter                              consul-exporter   running   running  false    service  10.40.11.85:29107(http)      ad77c86e  6fbd453a
pushgateway                                  pushgateway       running   running  false    service  10.40.11.85:29091(http)      5482dbb7  6fbd453a
web-express-26                               web-express-26    running   running  false    service  10.40.11.85:20746(http)      8ac9d894  6fbd453a
htdocs-26                                    htdocs-26         running   running  false    service  10.40.11.85:31231(http)      daf4d7d4  6fbd453a
redis-26                                     redis-26          running   running  false    service  10.40.11.85:25155(redis)     3e7f5150  6fbd453a
deploy-api-26                                deploy-api-26     running   running  false    service  10.40.11.85:25597(http)      5bc109fe  6fbd453a
user-postgres-26                             user-postgres-26  running   running  false    service  10.40.11.85:26287(postgres)  33fcf0f6  6fbd453a
nomad-exporter                               nomad-exporter    running   running  false    service  10.40.11.85:29172(http)      9d9959e3  6fbd453a
cadvisor                                     cadvisor          running   running  false    system   10.40.5.247:28080(http)      3ec5d326  2fc75e14
cadvisor                                     cadvisor          running   running  false    system   10.40.11.85:28080(http)      59766ca8  6fbd453a
pancakes                                     pancakes          running   running  false    service  10.40.11.85:23896(http)      bc126774  6fbd453a
pancakes                                     pancakes          running   running  false    service  10.40.11.85:30979(http)      aa5c6ec4  6fbd453a
ecr-exporter                                 ecr-exporter      running   running  false    service  10.40.11.85:28070(http)      bc89cb24  6fbd453a
```


Building & Installation
-----------------------

```
make
sudo make install
sudo make uninstall
```
