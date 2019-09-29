package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"math/rand"

	"github.com/garyburd/redigo/redis"
	"time"
)

type MainController struct {
	beego.Controller
}

const INF int =9999//最大路径长度
const MAX  int=99 //最大抽样数量

type Seat struct {
	No     int //座位号
	Name   string //财产的名称
	Amount int32  //财产的数额
}//座位及相关信息

type Range struct {
	FirstNo int//第一个座位编号
	SecondNo int//第二个座位编号
	Distance int//两点距离
}//两个座位间距离

/*此函数用于获取随机的抽样座位号，传入prop为抽样百分比，all为座位总数，返回存储有被选座位号的数组*/
func GetRandom(prop int) []int  {
	//计算目前mongoDB数据库中的固定资产数目
	session, err := mgo.Dial("mongodb://localhost")
	if err != nil {

		panic(err)

	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("test").C("Seats")
	var test []Seat
	err=c.Find(nil).All(&test)
	all:=len(test)-1//减去出发点
	count:=int(prop*all/100)
	start,end:=1,all
	if end < start || (end-start) < count {
		return nil
	}
	//存放结果的slice
	chosen := make([]int, 0)
	//随机数生成器，加入时间戳保证每次生成的随机数不一样
	rand.New(rand.NewSource(time.Now().UnixNano()))
	//r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(chosen) < count {
		//生成随机数
		num := rand.Intn(end - start) + start
		//查重
		exist := false
		for _, v := range chosen {
			if v == num {
				exist = true
				break
			}
		}

		if !exist {
			chosen = append(chosen, num)
		}
	}
	//将数组第一个置0
	fix :=make([]int ,count+1)
	for k:=0;k<count ;k++  {
		fix[k+1]=chosen[k]
	}
	fix[0]=0
	return fix
}

/*此函数用于生成最短路径*/
/*参数说明：dep:搜索的深度，level：层数，path:当前路径权值，matrix：邻接矩阵，bestPath:当前最优路径权值,bestPathX:当前最优路径,pathX:当前路径*/
func TSP(dep,level,path int, matrix [][MAX]int,bestPath *int,bestPathX,pathX []int){
	if dep == level {//搜索完成
		if path + matrix[pathX[level-1]][pathX[0]]<(*bestPath) {
			*bestPath=path + matrix[pathX[dep - 1]][pathX[0]]
			for j:=0;j<level ;j++  {
				bestPathX[j]=pathX[j]
			}
		}
	} else{//未搜索完成
		for j:=dep;j<level ;j++  {
			if path+matrix[pathX[dep-1]][pathX[j]] < (*bestPath) {
				bestPathX[dep],bestPathX[j]=bestPathX[j],bestPathX[dep]
				path += matrix[pathX[dep - 1]][pathX[dep]]
				TSP(dep+1,level,path, matrix,bestPath,bestPathX,pathX)
				path -= matrix[pathX[dep - 1]][pathX[dep]]
				pathX[dep],pathX[j]=pathX[j],pathX[dep]
			}
		}
	}
}

/*此函数用于按资产编号输出最短路径，matrix：邻接矩阵，chosen:被抽样的座位数组*/
func ShortestPathDisplay(matrix [][MAX]int,chosen []int,this *MainController){
	bestPath:=INF
	num:=len(chosen)
	bestPathX:=make([]int ,num)
	pathX:=make([]int,num)
	for j:=0;j<num ;j++  {
		bestPathX[j]=j
		pathX[j]=j
	}
	source,path:=0,0
	TSP(1,num,path, matrix,&bestPath,bestPathX,pathX)
	this.Data["bestPath"]=bestPath
	display:=make([]int,num+1)
	for i:=0;i<num ;i++  {
		display[i]=chosen[bestPathX[i]]
	}
	display[num]=source
	this.Data["pathDisplay"]=display
}

/*此函数用于生成各抽样点距离的邻接矩阵，chosen：被抽样的数组*/
func GetMatrix(chosen []int ) [][MAX]int {
	num:=len(chosen)
	matrix:=make([][MAX]int,num)
	for i:=0;i<num ;i++  {
		for j:=0;j<num ;j++  {
			if i==j {
				matrix[i][j]=0
			}else {
				FillMatrix(matrix,i,j,chosen)
			}
		}
	}
	return matrix
}

/*此函数用于将数据录入MongoDB*/
func SaveMongoDB(){
	session, err := mgo.Dial("mongodb://localhost")
	if err != nil {

		panic(err)

	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("Property").C("Seats")
		err = c.Insert(&Seat{
			No:     0,
			Name:   "Origin",
			Amount: 0,
		},&Seat{
			No:     1,
			Name:   "Property1",
			Amount: 1000,
		},&Seat{
			No:     2,
			Name:   "Property2",
			Amount: 2000,
		},&Seat{
			No:     3,
			Name:   "Property3",
			Amount: 4000,
		},&Seat{
			No:     4,
			Name:   "Property4",
			Amount: 5000,
		},&Seat{
			No:     5,
			Name:   "Property5",
			Amount: 4593,
		},&Seat{
			No:     6,
			Name:   "Property6",
			Amount: 7564,
		},&Seat{
			No:     7,
			Name:   "Property7",
			Amount: 5588,
		},&Seat{
			No:     8,
			Name:   "Property8",
			Amount: 8970,
		},&Seat{
			No:     9,
			Name:   "Property9",
			Amount: 4000,
		},&Seat{
			No:     10,
			Name:   "Property10",
			Amount: 2000,
		})
	c=session.DB("Property").C("Distance")
	rand.New(rand.NewSource(time.Now().UnixNano()))
	for i:=0;i<10 ;i++  {
		for j := 0; j < 10; j++ {
			if i!=j{
			err=c.Insert(&Range{
				FirstNo:  i,
				SecondNo: j,
				Distance:rand.Intn(99)+1 ,
			})
			}
		}
	}

}
/*此函数用于将查询redis数据，若不存在则在mongoDB数据库中查找，并填充数组，并写入redis缓存，redis缓存60秒过期*/
func FillMatrix(matrix [][MAX]int,i,j int,chosen []int) {
	conn, _ := redis.Dial("tcp", "127.0.0.1:6379")
	defer conn.Close()
	key:=string(i)+string(j)
	isExist,_:=redis.Bool(conn.Do("EXISTS",key))
	if isExist {
		//如果有
		// 从redis里直接读取
		result, _ := redis.Int(conn.Do("GET", key))
		matrix[i][j]= result
	}else {
		session, err := mgo.Dial("mongodb://localhost")
		if err != nil {
			panic(err)
		}
		defer session.Close()
		session.SetMode(mgo.Monotonic, true)
		c := session.DB("Property").C("Distance")
		result:=Range{}
		err = c.Find(bson.M{"firstno": chosen[i],"secondno":chosen[j]}).One(&result)
		matrix[i][j]=result.Distance
		//写入redis并且设置过期时间
		_, err = conn.Do("SET", key,result.Distance , "EX", "60")
		if err != nil {
			fmt.Println("redis set failed:", err)
		}
	}
}





func (this *MainController) Get() {
	SaveMongoDB()
	this.TplName = "index.html"
}
func (this *MainController) Post(){
	prop,_:=this.GetInt("prop")
	chosen:=GetRandom(prop)
	mat:= GetMatrix(chosen)
	ShortestPathDisplay(mat,chosen,this)
	this.TplName = "display.html"
}

