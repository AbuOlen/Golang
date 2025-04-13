package main

import (
	"fmt"
	users "hw5/users"
)

func main() {
	srv := users.NewService()

	usr1, err := srv.CreateUser("Nazar")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(usr1)
	usr2, err := srv.CreateUser("Marko")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(usr2)
	usr3, err := srv.CreateUser("Zenyk")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(usr3)

	fmt.Println(srv.GetUser(usr1.ID))

	fmt.Println(srv.DeleteUser(usr2.ID))

	fmt.Println(srv.ListUsers())
}

//&{f94c5bbb-bb4a-4fab-abda-eab3e81bcea4 Nazar}
//&{e7863464-ed02-4301-bb69-d896277cf4e6 Marko}
//&{5432b9b0-5d26-4989-8aa3-1783daadd255 Zenyk}
//&{f94c5bbb-bb4a-4fab-abda-eab3e81bcea4 Nazar} <nil>
//<nil>
//[{5432b9b0-5d26-4989-8aa3-1783daadd255 Zenyk} {f94c5bbb-bb4a-4fab-abda-eab3e81bcea4 Nazar}] <nil>
