# Dolpin 

> still development

Dolpin is auto generate test case and very simple to use it
example we have code like this
``` go

func IsUserExist(userId int) bool {
    // get data to db
    user := userrepository.New()
    user, err := user.GetUserByID(userId)
    if err != nil || user == nil {
        return false
    }
    return true
}

```
instead we create unit test like this
``` go
func TestIsUserExist(t *testing.T){

    var userId = 123

    // setup mock
    mock := userReposiotyMock.New()
    // expecting call fucntion
    
    type Arg struct {
        userId int
    }
    var testCase = []struct{
        Name string
        Arg Arg
        Prepare func(t *testing.T, arg Arg) {}
        Expected bool
    }{
        {
            Name: "should be success",
            Arg: Arg{
                userId: userId,
            },
            Prepare: func(t *testing.T, arg Arg){
                mock.Expect("GetUserById").WithArg(arg.UserId).Return(mock.Anything, nil)
            },
            Expect: true,
        },
        {
            Name: "should be error on function GetUserById",
            Arg: Arg{
                userId: userId,
            },
            Prepare: func(t *testing.T, arg Arg){
                mock.Expect("GetUserById").WithArg(arg.UserId).Return(mock.Anything, errors.New("something error"))
            },
            Expect: false,
        },
        {
            Name: "should be return nil in fild user on GetUserById",
            Arg: Arg{
                userId: userId,
            },
            Prepare: func(t *testing.T, arg Arg){
                mock.Expect("GetUserById").WithArg(arg.UserId).Return(nil, nil)
            },
            Expect: false,
        }
    }

    for _, tc := range testCase {

        t.Run(tc.Name, func(t *testing.T){
            tc.Prepare(t, tc.Arg)

            isExist := IsUserExist(tc.Arg.userId)

            if isExist != tc.Expect {
                t.Fatal("expect %v, got %v", tc.Expect, isExist)
            }
        })
    }
}
```

with dolpin we just create
``` go
func TestIsUserExist(t *testing.T){

    _dolpin := dolpin.New(dolpin.Config{
        MockFolder: "/mocks"
    })
    // prepare function is for calling moc
    _dolpin.CallFunction("GetUserById", dolpin.Args(dolpin.Anything), dolpin.Results(nil))

    // prepare arg and expect value 
    _dolpin.Prepare(dolpin.Arg(123), dolpin.Expect(true))
    _dolpin.Prepare(dolpin.Arg(456), dolpin.Expect(false))
    _dolpin.SetDecoration("type IsUserExist func(int) bool")
    _dolpin.TestTarget(t,IsUserExist)
}
```