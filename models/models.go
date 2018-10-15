package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Id       int `orm:"pk;auto"`
	Username string
	Password string
	Tel      string
	Balance  float64
}
type Product struct {
	Id      int `orm:"pk;auto"`
	Name    string
	Stock   int     //库存数量
	Price   float64 `orm:"digits(8);decimals(2)"`
	Remarks string
	Img     string
}
type Item struct {
	Id      int `orm:"pk;auto"`
	Pid     int //关联产品ID
	Uid     int //关联使用者ID
	Oid     int //关联订单ID
	Remarks string
	Flag    int
}

type Itemnum struct {
	Id     int `orm:"pk;auto"`
	Itemid int
	Num    int
}

type Order struct {
	Id         int     `orm:"pk;auto"`
	Uid        int     //关联使用者ID
	Totalprice float64 `orm:"digits(8);decimals(2)"`
	Remarks    string
}

func RegisterDB() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "root:123456@/shopping?charset=utf8&loc=Local")
	orm.RegisterModel(new(User), new(Product), new(Item), new(Itemnum), new(Order))
}

func (c *Product) AddProduct(name string, price float64, stock int, remark, img string) error {
	o := orm.NewOrm()
	product := &Product{
		Name:    name,
		Price:   price,
		Stock:   stock,
		Remarks: remark,
		Img:     img,
	}
	o.QueryTable("product").Filter("name",name).One(product)
	_, err := o.Insert(product)
	return err
}

func GetAllProduct() ([]*Product, error) {
	o := orm.NewOrm()
	products := make([]*Product, 0)
	qs := o.QueryTable("product")
	var err error
	_, err = qs.All(&products)
	return products, err
}

func GetAllOrders(uid int) ([]*Order, error) {
	o := orm.NewOrm()
	orders := make([]*Order, 0)
	qs := o.QueryTable("order")
	var err error
	_, err = qs.Filter("uid", uid).All(&orders)
	return orders, err
}

//获取用户信息
func GetUserInfo(username string) *User {
	o := orm.NewOrm()
	user := new(User)
	qs := o.QueryTable("user")
	qs.Filter("username", username).One(user)
	return user
}

//添加用户
func (u *User) Add(username string, password string, telephone string) (int64, error) {

	o := orm.NewOrm()
	var user User
	user.Username = username
	user.Password = password
	user.Tel = telephone
	id, err := o.Insert(&user)
	return id, err
}

//新增item，如果pid，uid，flag=0则更新
func (i *Item) Additem(pid int, uid int, num int,remarks string) string {
	o := orm.NewOrm()
	var item Item
	err := o.QueryTable("Item").Filter("pid", pid).Filter("uid", uid).Filter("flag", int(0)).One(&item)
	// 多条的时候报错
	if err == orm.ErrMultiRows {

		return "多条数据错误"
	}
	// 没有找到记录
	if err == orm.ErrNoRows {
		item.Pid = pid
		item.Remarks = remarks
		item.Flag = int(0)
		item.Uid = uid
		item.Oid = int(0)
		id, ierr := o.Insert(&item)
		if ierr == nil {
			var itns Itemnum
			_, itnerr := itns.Additem(int(id), num)
			if itnerr == nil {
				return "添加成功"
			}
			return "数据插入错误"
		}
		return "数据插入错误"

	}
	res, errrow := o.Raw("UPDATE Itemnum SET num = num + ? WHERE Itemid = ?", num, item.Id).Exec()
	if errrow == nil {
		res.RowsAffected()
		return "更新成功"
	}

	return "更新失败"
}

func (n *Itemnum) Addmin(itid int, num int) string {
	o := orm.NewOrm()
	res, errrow := o.Raw("UPDATE Itemnum SET num = ? WHERE Itemid = ?", num, itid).Exec()
	if errrow == nil {
		res.RowsAffected()
		return "更新成功"
	}
	return "更新失败"
}

func (n *Order) Addorder(ids []string, uid int, remarks string, totalp float64) string {
	o := orm.NewOrm()
	var order Order
	order.Uid = uid
	order.Totalprice = totalp
	order.Remarks = remarks
	//var idsarry = strings.Split(ids, ",")
	//var aaa string
	var str = "UPDATE Item SET flag = 2,oid=? WHERE Id in (0"
	for _, v := range ids {
		str = str + "," + v
	}
	str = str + ")"
	id, err := o.Insert(&order)
	if err == nil {
		res, errrow := o.Raw(str, id).Exec()
		if errrow == nil {
			res.RowsAffected()
			return "更新数据成功"
		}
		return "更新数据失败"
	}
	return "更新数据失败"

}

func (n *Item) Deleteitem(itid int) string {
	o := orm.NewOrm()
	res, errrow := o.Raw("UPDATE Item SET flag = ? WHERE Id = ?", 1, itid).Exec()
	if errrow == nil {
		res.RowsAffected()
		return "更新成功"
	}
	return "更新失败"

}

func (n *Itemnum) Additem(itid int, num int) (int64, error) {
	o := orm.NewOrm()
	var itn Itemnum
	itn.Itemid = itid
	itn.Num = num
	id, err := o.Insert(&itn)
	return id, err
}

func Queryitems(uid int) []orm.Params {
	o := orm.NewOrm()
	var maps []orm.Params
	o.Raw("select it.id as itemid,pro.`name` as proname,pro.img as img,it.remarks as remarks,itn.num as num ,pro.price as price,itn.num * pro.price as total from item as it left join itemnum as itn on it.id = itn.itemid left join product as pro on it.pid = pro.id where it.flag =? and it.uid = ?", 0, uid).Values(&maps)
	return maps
}

func UpdateTheBalance(balance, total float64, uid int) {
	o := orm.NewOrm()
	user := User{Id: uid}
	user.Balance = balance - total
	o.Update(&user, "Balance")
}

func GetBalanceById(id int) float64 {
	o := orm.NewOrm()
	user := new(User)
	o.QueryTable("user").Filter("id", id).One(user)
	return user.Balance
}
