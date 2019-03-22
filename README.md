Run command on multiple servers. Thanks to this little [snippet] from [@Erik](https://github.com/erikdubbelboer/).

```
go get github.com/xoxzo/myssh/cmd/myssh
$GOPATH/bin/myssh -h host-1,host-2
```

## With gore
Use [gore].

```
go get github.com/xoxzo/myssh
$GOPATH/bin/gore
gore> :import github.com/xoxzo/myssh
gore> var hosts []string
gore> hosts = append(hosts, "host-1:22")
[]string{"host-1:22"}
gore> hosts = append(hosts, "host-2:22")
[]string{"host-1:22", "host-2:22"}
gore> myssh.Run("hostname", hosts)
Running hostname at host-1:22
host-1:22: host-1

Running hostname at host-2:22
host-2:22: host-2

gore>
```

[gore]:https://github.com/motemen/gore
[snippet]:https://gist.github.com/erikdubbelboer/f62a109d8e8798a11eb89ed494491953
