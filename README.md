# gossh
golang实现的一个ssh远程执行小工具

- 依赖：
  ```
  require (
    github.com/pkg/sftp v1.13.5 // sftp连接工具包
    golang.org/x/crypto v0.11.0 // ssh连接工具包
  )
  ```

# 使用示例
1. 使用username+password连接到远端服务器，并执行shell
  ```go
  func TestGetRemote1(t *testing.T) {
    con, err := gossh.Remote1("root", "password", "远端服务器ip:22")
    if err != nil {
      t.Error(err)
      return
    }
    defer con.Close()
    s, err2 := con.ExecShell(context.Background(), "df -h")
    if err2 != nil {
      t.Error(err2)
      return
    }
    t.Logf(s)
  }
  ```
2. 使用username+私钥连接到远端服务器，并执行shell
  ```go
  func TestGetRemote2(t *testing.T) {
    con, err := gossh.Remote2("root", "/root/.ssh/id_rsa", "xxx.xxx.xxx.xxx:22")
    if err != nil {
      t.Error(err)
      return
    }
    defer con.Close()
    s, err2 := con.ExecShell(context.Background(), "df -h")
    if err2 != nil {
      t.Error(err2)
      return
    }
    t.Logf(s)
  }
  ```
3. 本地执行shell
  ```go
  func TestGetLocal(t *testing.T) {
    l := gossh.Local()
    defer l.Close()
    s, err := l.ExecShell(context.Background(), "df -h")
    if err != nil {
      t.Error(err)
      return
    }
    t.Logf(s)
  }
  ```
4. 使用默认方式远程连接到远端服务器，并执行shell
  ```go
  func TestGetRemoteDefault(t *testing.T) {
    con, err := gossh.RemoteDefault("目标机器ip:22")
    if err != nil {
      t.Error(err)
      return
    }
    defer con.Close()
    s, err2 := con.ExecShell(context.Background(), "df -h")
    if err2 != nil {
      t.Error(err2)
      return
    }
    t.Logf(s)
  }
  ```
  **注意**：该方式需要将下面`公钥`添加到目标机器的`/root/.ssh/authorized_keys`文件中
  ```
  ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDP3qG9zzNQkhKuWWJTCHAuY04bQ9h/vhZplVrTSnEoWL1SsbT7v/dCDXuyazNDo9ikd9BS6H/nE5lOKp+Omi1W2uWs/kxrCNhouXRO8kMLViRB3DRP2VYFDo36UafzNdGKkH/vW4Ptilga/ucForW05SindT3KeKf+tB1u3RBlRz6rzpeuqrflVtcaWtQ33exWMO8CxgzCtsDexWWLP+TLdeOaWyfn0hj4tf36+K7oENAzGGhQuEwETiMUkJKfykBThBenWgU9mM1/5VbgvGiW7xIoeyDX8RI6Lz5q8mb3+ajuEqPyX/qwiNasYkQ7bWGaLDAVF3yJ20w7EpP54yi9rEoiBt6GAEo2JX5OuibzwMsz2CCykiB8H4YyiOlBY0q5GwrXC87fslvEt4KcYdk/XrZT+ikrJePgTCbQJhUGf8yYe1aDKoTBf/uIuT/O7aJ49KxnfeQdxel3xIykyROxnNisQ8Iz3vdC/QZlZsQUnJzo0UXtwDpwmwLwjpLFCMM= gossh@github.com
  ```