#smtp
##smtp

    package main

    import (
        "log"

        "github.com/ProForks/toolkits/smtp"
    

    func main() {
        auth := smtp.LoginAuth("notify@a.com", "password","smtp.exmail.qq.com:25")
        log.Println(auth.SendMail("notify@a.com", "ulric@b.com;rain@c.com", "这是subject", "这是body,<font color=red>red</font>"))
    }
