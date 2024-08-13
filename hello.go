package main

import "fmt"


var users = [3] string {"sdfw", "dwfwefwf", "333333"}
var n = [10] int {3,4,5,2,3,2,4}

func main() {


users := []string{"Bob", "Alice", "Kate", "Sam", "Tom", "Paul", "Mike", "Robert"}
//удаляем 4-й элемент
var n = 3
users = append(users[:n], users[n+1:]...)
fmt.Println(users)
  
}


func summ(numbers ...int) int{
    nout := 0
    for _,n := range numbers{
        nout += n
    }
    return nout
}