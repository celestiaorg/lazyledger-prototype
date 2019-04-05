package lazyledger

import (
    "testing"
)


func TestAppDummySimpleBlock(t *testing.T) {
    bs := NewSimpleBlockStore()
    b := NewBlockchain(bs)

    sb := NewSimpleBlock([]byte{0})

    ms := NewSimpleMap()
    app := NewDummyApp(ms)
    var regApp Application
    regApp = app
    b.RegisterApplication(&regApp)

    puts := make(map[string]string)
    puts["foo"] = "bar"
    puts["goo"] = "tar"

    sb.AddMessage(app.GenerateTransaction(puts))
    b.ProcessBlock(sb)

    if app.Get("foo") != "bar" || app.Get("goo") != "tar" {
        t.Error("dummy app state update failed")
    }
}
