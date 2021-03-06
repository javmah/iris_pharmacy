package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/kataras/iris"

	"database/sql"

	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
	"github.com/kataras/iris/sessions"

	_ "github.com/go-sql-driver/mysql"
	// Jwt
	//"github.com/dgrijalva/jwt-go"
)

var (
	cookieNameForSessionID = "mycookiesessionnameid"
	sess                   = sessions.New(sessions.Config{Cookie: cookieNameForSessionID})
)

type adduserdata struct {
	FirstName   string
	LastName    string
	MiddleName  string
	Gender      string
	Dob         string
	Designation string
	Area        string
	Postcode    string
	Username    string
	Password    string
}

type showUserdata struct {
	ID          string
	FirstName   string
	LastName    string
	MiddleName  string
	Gender      string
	Dob         string
	Designation string
	Area        string
	Postcode    string
	Username    string
	Password    string
}

type addproduct struct {
	Tradenames    string
	Genericnames  string
	Chemicalnames string
	Activationsta string
	UsedFor       string
	Mrp           string
}

type products struct {
	Id            string
	Tradenames    string
	Genericnames  string
	Chemicalnames string
	Activationsta string
	UsedFor       string
	Mrp           string
}

func main() {
	app := iris.New()
	app.Logger().SetLevel("debug")
	// Optionally, add two built'n handlers
	// that can recover from any http-relative panics
	// and log the requests to the terminal.
	app.Use(recover.New())
	app.Use(logger.New())

	// javmah Starts
	app.RegisterView(iris.HTML("./templates", ".html"))
	app.StaticWeb("/", "./assets")
	// iris.StaticServe("./static", "/static").Binary(Asset, AssetNames)

	// 1. Routes
	app.Get("/", func(ctx iris.Context) {
		ctx.HTML(" <a href='signin'> Sign In  </a> <br> <a href='dashboard'> Dashboard  </a> <br> <a href='users'> user </a>  <br> <a href='add-user'> add user  </a>  <br> <a href='products'> product   </a>  <br> <a href='add-product'> Add Product   </a>  ")
	})

	// 2. Routes
	app.Get("/signin", func(ctx iris.Context) {

		ctx.View("login.html")
	})

	// 3. Login Validation
	app.Post("/signinval", func(ctx iris.Context) {

		Username := ctx.PostValue("Username")
		Password := ctx.PostValue("Password")

		var id string
		var First string
		var Last string
		var Desig string
		var area string
		var Postc string

		if Username != "" && Password != "" {
			db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/medicine")
			if err != nil {
				fmt.Println("Error On Connection :", err.Error())
			}

			var numberOfrowEffected int
			rows, err := db.Prepare("SELECT COUNT(*) as count FROM user WHERE (Username=? AND Password=?) ")
			rows.QueryRow(Username, Password).Scan(&numberOfrowEffected)

			if err != nil {
				println("ERROR : ", rows)
				log.Fatal(err)
			}

			if numberOfrowEffected == 1 {
				userdetails, erruserdata := db.Prepare("SELECT user.id, user.FirstName,  user.LastName, user.Designation,  user.Area,  user.Postcode FROM  user  WHERE Username=? AND Password=?  ")
				if erruserdata != nil {
					fmt.Println("ERROR Gatting User Info : ", erruserdata)
				}
				userdetails.QueryRow(Username, Password).Scan(&id, &First, &Last, &Desig, &area, &Postc)
				// Session Starts
				session := sess.Start(ctx)
				// session.Set("authenticated", "1")
				session.Set("authenticated", true)

				session.Set("id", id)
				session.Set("Firstname", First)
				session.Set("Lastname", Last)
				session.Set("Desig", Desig)
				session.Set("Area", area)
				session.Set("Postc", Postc)
			}

			fmt.Println("Number of Query Row Effected : ", numberOfrowEffected)
			//  Redirect if Logged In
			ctx.Redirect("/dashboard")

		} else {
			// ctx.Writef("Empty || Empty ")
			ctx.Redirect("/signin")
		}

	})

	// 3.  Dashboard
	app.Get("/dashboard", func(ctx iris.Context) {
		// Cheack Valid User Starts
		// Check if user is authenticated
		if auth, _ := sess.Start(ctx).GetBoolean("authenticated"); !auth {
			ctx.Redirect("/signin")
			return
		}
		// Cheack Valid User Starts

		ctx.View("index.html")
	})

	// 4. For User
	app.Get("/users", func(ctx iris.Context) {
		// if auth, _ := sess.Start(ctx).GetBoolean("authenticated"); !auth {
		// 	ctx.Redirect("/signin")
		// 	return
		// }

		// Databsae
		db, dbconnerr := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/medicine")
		// SELECT user.id, user.FirstName,  user.LastName, user.Designation,  user.Area,  user.Postcode FROM  user  WHERE Username=? AND Password=?
		rows, err := db.Query("SELECT user.id, user.FirstName,  user.LastName, user.MiddleName, user.Gender , user.Dob, user.Designation, user.Area, user.Postcode ,user.Username FROM  user ")
		// rows, err := db.Query("SELECT * FROM user ")

		if dbconnerr != nil {
			fmt.Println("Hmm err :", err.Error())
		}
		var result []showUserdata
		for rows.Next() {
			student := showUserdata{}
			err2 := rows.Scan(&student.ID, &student.FirstName, &student.LastName, &student.MiddleName, &student.Gender, &student.Dob, &student.Designation, &student.Area, &student.Postcode, &student.Username)
			if err2 != nil {
				panic(err2)
			}
			result = append(result, student)
		}
		fmt.Println(result)
		ctx.ViewData("result", result)

		ctx.View("users.html")
	})

	// 5. For ADD a User
	app.Get("/add-user", func(ctx iris.Context) {
		// Check if user is authenticated Starts
		// if auth, _ := sess.Start(ctx).GetBoolean("authenticated"); !auth {
		// 	ctx.Redirect("/signin")
		// 	return
		// }
		// Check if user is authenticated Ends
		ctx.View("add-user.html")
	})
	// POST Form Process
	// For Validation Help
	// https://github.com/go-playground/validator
	// https://gopkg.in/bluesuncorp/validator.v9
	app.Post("/saveadduserdata", func(ctx iris.Context) {
		userid := ctx.PostValue("ID")
		if userid != "" {
			data := adduserdata{}
			err := ctx.ReadForm(&data)
			if err != nil {
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.WriteString(err.Error())
			}
			// inserting user to database
			db, dbconnerr := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/medicine")
			if dbconnerr != nil {
				fmt.Println("Hmm err :", err.Error())
			} else {
				update, qperr := db.Prepare("UPDATE  user SET FirstName =? , LastName =? , MiddleName =? , Gender =? , Dob =? ,Designation =? ,Area =?,Postcode =?,Username =?, Password =? WHERE id= ?")
				if qperr != nil {
					panic(qperr.Error())
				}

				// hashe Passwoard starts
				passw := data.Password
				h := md5.New()
				h.Write([]byte(passw))
				hashedPass := hex.EncodeToString(h.Sum(nil))
				fmt.Println("MD5 hases Are :", hashedPass)
				// hashe Passwoard Ends

				_, exerr := update.Exec(
					data.FirstName,
					data.LastName,
					data.MiddleName,
					data.Gender,
					data.Dob,
					data.Designation,
					data.Area,
					data.Postcode,
					data.Username,
					hashedPass,
					userid,
				)

				if exerr != nil {
					panic(exerr.Error())
				}

			}
			defer db.Close()
		} else {
			data := adduserdata{}
			err := ctx.ReadForm(&data)
			if err != nil {
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.WriteString(err.Error())
			}
			// inserting user to database
			db, dbconnerr := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/medicine")
			if dbconnerr != nil {
				fmt.Println("Hmm err :", err.Error())
			} else {
				stmtIns, qperr := db.Prepare("INSERT INTO user SET FirstName =? , LastName =? , MiddleName =? , Gender =?, Dob =? ,Designation =? ,Area =?,Postcode =?,Username =?, Password =?")
				if qperr != nil {
					panic(qperr.Error())
				}

				// hashe Passwoard starts
				passw := data.Password
				h := md5.New()
				h.Write([]byte(passw))
				hashedPass := hex.EncodeToString(h.Sum(nil))
				fmt.Println("MD5 hases Are :", hashedPass)
				// hashe Passwoard Ends

				_, exerr := stmtIns.Exec(
					data.FirstName,
					data.LastName,
					data.MiddleName,
					data.Gender,
					data.Dob,
					data.Designation,
					data.Area,
					data.Postcode,
					data.Username,
					hashedPass,
				)

				if exerr != nil {
					panic(exerr.Error())
				}

			}
			defer db.Close()
		}

		// ctx.Writef("Visitor: %#v", data)

	})

	// 6. Edit User
	app.Get("/edit-user/{id:int}", func(ctx iris.Context) {
		// if auth, _ := sess.Start(ctx).GetBoolean("authenticated"); !auth {
		// 	ctx.Redirect("/signin")
		// 	return
		// }

		userinfo := showUserdata{}
		userID, _ := ctx.Params().GetInt("id")
		// Databsae
		db, _ := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/medicine")
		row, _ := db.Prepare("SELECT user.id, user.FirstName, user.LastName, user.MiddleName, user.Gender , user.Dob, user.Designation, user.Area, user.Postcode ,user.Username ,user.Password FROM  user WHERE user.id = ?")
		row.QueryRow(userID).Scan(&userinfo.ID, &userinfo.FirstName, &userinfo.LastName, &userinfo.MiddleName, &userinfo.Gender, &userinfo.Dob, &userinfo.Designation, &userinfo.Area, &userinfo.Postcode, &userinfo.Username, &userinfo.Password)

		// fmt.Println(userinfo)
		ctx.ViewData("", userinfo)
		ctx.View("edit-user.html")
	})

	// 6. Product
	// https://stackoverflow.com/questions/23304854/how-do-you-determine-if-a-variable-is-a-slice-or-array
	app.Get("/products", func(ctx iris.Context) {
		// Check if user is authenticated Starts
		// if auth, _ := sess.Start(ctx).GetBoolean("authenticated"); !auth {
		// 	ctx.Redirect("/signin")
		// 	return
		// }
		// Databsae
		db, dbconnerr := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/medicine")
		rows, err := db.Query("SELECT * FROM medicinelist")
		if dbconnerr != nil {
			fmt.Println("Hmm err :", err.Error())
		}
		var result []products
		for rows.Next() {
			item := products{}
			err2 := rows.Scan(&item.Id, &item.Tradenames, &item.Genericnames, &item.Chemicalnames, &item.Activationsta, &item.UsedFor, &item.Mrp)
			if err2 != nil {
				panic(err2)
			}
			result = append(result, item)
		}
		fmt.Println(result)
		ctx.ViewData("result", result)
		ctx.View("products.html")
	})

	// 7. Add-product
	app.Get("/add-product", func(ctx iris.Context) {
		// Check if user is authenticated Starts
		// if auth, _ := sess.Start(ctx).GetBoolean("authenticated"); !auth {
		// 	ctx.Redirect("/signin")
		// 	return
		// }
		// Check if user is authenticated Ends
		ctx.View("add-product.html")
	})

	// 8. edit-product
	app.Get("/edit-product/{id:int}", func(ctx iris.Context) {
		// Check if user is authenticated Starts
		// if auth, _ := sess.Start(ctx).GetBoolean("authenticated"); !auth {
		// 	ctx.Redirect("/signin")
		// 	return
		// }
		// Check if user is authenticated Ends

		productId, _ := ctx.Params().GetInt("id")
		editProduct := products{}
		// Databsae
		db, _ := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/medicine")
		row, _ := db.Prepare("SELECT medicinelist.id ,  medicinelist.TradeNames, medicinelist.GenericNames, medicinelist.ChemicalNames, medicinelist.ActivationStatus , medicinelist.UsedFor, medicinelist.Mrp  FROM  medicinelist WHERE medicinelist.id = ?")
		row.QueryRow(productId).Scan(&editProduct.Id, &editProduct.Tradenames, &editProduct.Genericnames, &editProduct.Chemicalnames, &editProduct.Activationsta, &editProduct.UsedFor, &editProduct.Mrp)

		// fmt.Println(userinfo)
		fmt.Println(editProduct)
		ctx.ViewData("", editProduct)
		ctx.View("edit-product.html")
	})

	// 9. Orders
	// app.Get("/orders", func(ctx iris.Context) {
	// 	// Check if user is authenticated Starts
	// 	// if auth, _ := sess.Start(ctx).GetBoolean("authenticated"); !auth {
	// 	// 	ctx.Redirect("/signin")
	// 	// 	return
	// 	// }
	// 	// Check if user is authenticated Ends
	// 	ctx.View("orders.html")
	// })

	app.Get("/orders", func(ctx iris.Context) {

		type product struct {
			OrderId    int
			Productid  int
			TradeNames string
			Qty        int
			Price      int
		}
		type orders struct {
			Id           int
			UserId       int
			FirstName    string
			OrderDate    string
			DeliveryDate string
			Products     []product
		}

		// Databsae
		db, dbconnerr := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/medicine")
		if dbconnerr != nil {
			fmt.Println("Hmm err :", dbconnerr)
		}

		rows, err := db.Query("SELECT orders.Id , orders.UserId, user.FirstName , orders.OrderDate , orders.DeliveryDate FROM orders INNER JOIN user ON orders.UserId = user.id ORDER BY orders.Id ASC")
		if err != nil {
			fmt.Println("Hmm err :", err)
		}
		// Start faching Data
		var Results []orders
		for rows.Next() {
			order := orders{}
			err2 := rows.Scan(&order.Id, &order.UserId, &order.FirstName, &order.OrderDate, &order.DeliveryDate)
			if err2 != nil {
				panic(err2)
			}
			// Gatting Order Items
			itemrows, erritems := db.Query("SELECT order_items.OrderId , order_items.Productid , medicinelist.TradeNames , order_items.Qty, order_items.Price   FROM order_items INNER JOIN medicinelist ON order_items.Productid = medicinelist.id WHERE order_items.OrderId = ? ORDER BY order_items.Productid ASC", order.Id)
			if erritems != nil {
				println("ERR: Error Gatting Item  ")
			}
			for itemrows.Next() {
				Items := product{}
				itemrows.Scan(
					&Items.OrderId,
					&Items.Productid,
					&Items.TradeNames,
					&Items.Qty,
					&Items.Price,
				)
				order.Products = append(order.Products, Items)
				// fmt.Println(Items)
			}
			Results = append(Results, order)
		}
		ctx.ViewData("result", Results)
		ctx.View("orders.html")
	})

	// 10. create-order
	app.Get("/create-order", func(ctx iris.Context) {
		// Check if user is authenticated Starts
		// if auth, _ := sess.Start(ctx).GetBoolean("authenticated"); !auth {
		// 	ctx.Redirect("/signin")
		// 	return
		// }
		// Check if user is authenticated Ends

		db, dbconnerr := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/medicine")
		rows, err := db.Query("SELECT * FROM medicinelist")
		if dbconnerr != nil {
			fmt.Println("Hmm err :", err.Error())
		}
		var result []products
		for rows.Next() {
			item := products{}
			err2 := rows.Scan(&item.Id, &item.Tradenames, &item.Genericnames, &item.Chemicalnames, &item.Activationsta, &item.UsedFor, &item.Mrp)
			if err2 != nil {
				panic(err2)
			}
			result = append(result, item)
		}
		fmt.Println(result)

		ctx.ViewData("result", result)
		ctx.View("create-order.html")
	})

	// ************************************************************************************************

	app.Post("/save-orders", func(ctx iris.Context) {
		db, conerr := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/medicine")
		if conerr != nil {
			fmt.Println(" SQL connection error :", conerr)
		}

		formValues := ctx.FormValues()

		// Deleting Empty Map  // Where Value Is Empty
		for k, _ := range formValues {
			if formValues[k][0] == "" {
				delete(formValues, k)
			}
		}

		// Inser Order table // create a New Order
		cretorder, creerr := db.Prepare("INSERT INTO orders(UserId) VALUES ( ? ) ")
		if creerr != nil {
			fmt.Println(" SQL Prepration error :", creerr)
		}
		// User id is
		sessiondataid := sess.Start(ctx).GetString("id")
		res, creorderexerr := cretorder.Exec(sessiondataid)
		if creorderexerr != nil {
			fmt.Println(" Execution error :", creorderexerr)
		}
		orderid, _ := res.LastInsertId()

		fmt.Println("Last Inserted Id is :", orderid)

		// Inserting Data Into order item table
		stmt, sqlpreperr := db.Prepare("INSERT INTO order_items( OrderId, Productid, Qty) VALUES ( ? , ? , ?  ) ")
		if sqlpreperr != nil {
			fmt.Println(" SQL Prepration error :", sqlpreperr)
		}
		//
		for k, _ := range formValues {
			res, excerr := stmt.Exec(orderid, k, formValues[k][0])
			if excerr != nil {
				fmt.Println(" Execution error :", excerr)
			}
			fmt.Println(res)
		}

		// fmt.Println(formValues)
		// fmt.Println(formValues["vehicle"][0])
		// fmt.Println(formValues["vehicle"][1])
		// formValues["vehicle"][1] = "BUS"
		// fmt.Println(formValues)
		// fmt.Println("_____________")
		// fmt.Println(formValues1)

	})

	//10. inventory
	app.Get("/inventory", func(ctx iris.Context) {
		// Check if user is authenticated Starts
		// if auth, _ := sess.Start(ctx).GetBoolean("authenticated"); !auth {
		// 	ctx.Redirect("/signin")
		// 	return
		// }
		// Check if user is authenticated Ends
		ctx.View("inventory.html")
	})

	//11. Reports
	app.Get("/reports", func(ctx iris.Context) {
		type qty struct {
			Id         int
			TradeNames string
			QTY        int
			TotalPrice int
		}
		// Databsae
		db, dbconnerr := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/medicine")
		if dbconnerr != nil {
			fmt.Println("ERR:DB Connection error : ", dbconnerr)
		}
		//
		itemrows, erritems := db.Query("SELECT medicinelist.id ,medicinelist.TradeNames , SUM(order_items.Qty) AS QTY , SUM(order_items.Price) as TotalPrice FROM medicinelist , order_items ,orders WHERE medicinelist.id = order_items.Productid AND orders.Id = order_items.OrderId AND orders.OrderDate = '2018-04-17 20:08:35' GROUP BY medicinelist.id")

		if erritems != nil {
			println("ERR : error While Gatting Product :", erritems)
		}

		var results []qty

		for itemrows.Next() {
			Items := qty{}
			itemrows.Scan(
				&Items.Id,
				&Items.TradeNames,
				&Items.QTY,
				&Items.TotalPrice,
			)
			results = append(results, Items)
		}

		fmt.Println(results)
		ctx.ViewData("result", results)
		ctx.View("reports.html")
	})

	//12. Settings
	app.Get("/settings", func(ctx iris.Context) {
		// Check if user is authenticated Starts
		// if auth, _ := sess.Start(ctx).GetBoolean("authenticated"); auth {
		// 	ctx.Redirect("/signin")
		// 	return
		// }
		// Check if user is authenticated Ends
		ctx.View("settings.html")
	})

	// 13. 404
	app.OnErrorCode(iris.StatusNotFound, func(ctx iris.Context) {
		// when 404 then render the template $views_dir/errors/404.html
		ctx.View("page_404.html")
	})

	// app.Get("/test/{id:int}", func(ctx iris.Context) {

	// 	userID, _ := ctx.Params().GetInt("id")
	// 	fmt.Println("USER id is :", userID)
	// })

	// ***************************  Test Starts *************************

	app.Get("/test", func(ctx iris.Context) {
		ctx.WriteString("  I have No idea What I am doing  ")
	})

	app.Handle(iris.MethodGet, "/test2", handler)

	// ***************************  Test Ends *************************

	// Add Product Form Process
	app.Post("/addproduct", func(ctx iris.Context) {
		addproductdat := addproduct{}
		ctx.ReadForm(&addproductdat)
		// database Starts
		db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/medicine")
		if err != nil {
			fmt.Println("Error On Connection :", err.Error())
		}
		stmtIns, qperr := db.Prepare("INSERT INTO 	medicinelist SET TradeNames =? , GenericNames =? , ChemicalNames =? ,ActivationStatus =?, UsedFor =? ,Mrp =? ")
		if qperr != nil {
			panic(qperr.Error())
		}
		_, exerr := stmtIns.Exec(
			addproductdat.Tradenames,
			addproductdat.Genericnames,
			addproductdat.Chemicalnames,
			addproductdat.Activationsta,
			addproductdat.UsedFor,
			addproductdat.Mrp,
		)
		if exerr != nil {
			panic(exerr.Error())
		}
		ctx.WriteString(" Hello from Add Product : 786 ")
	})

	// SignOut
	app.Get("/signout", func(ctx iris.Context) {
		//destroy, removes the entire session data and cookie
		sess.Destroy(ctx)
		ctx.Redirect("/signin")
	})

	// Delete
	app.Get("/delete/{rfrom:string}/{id:int}", func(ctx iris.Context) {
		db, _ := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/medicine")
		from := ctx.Params().Get("rfrom")
		id, _ := ctx.Params().GetInt("id")
		fmt.Println(" From :", from)
		fmt.Println(" id is :", id)
		// _____________________________

		switch from {
		case "users":
			stmt, _ := db.Prepare("delete from user where id=?")
			res, _ := stmt.Exec(id)
			affect, _ := res.RowsAffected()
			fmt.Println(affect)
			ctx.Redirect("/users")
		case "products":
			stmt, _ := db.Prepare("delete from medicinelist where id=?")
			res, _ := stmt.Exec(id)
			affect, _ := res.RowsAffected()
			fmt.Println(affect)
			ctx.Redirect("/products")
		default:
			fmt.Println("Err :  Nothin Is Delited , Altho it's called ")
		}

		// ______________________________

		// stmt, _ := db.Prepare("delete from user where id=?")
		// res, _ := stmt.Exec(userID)
		// affect, _ := res.RowsAffected()
		// ctx.Redirect("/users")

	})

	// API Srarts From heare ..................................................................................................................
	app.Get("/api/list", ApiList)
	app.Post("/api/order", ApiOrder)
	// API ENDs From heare ....................................................................................................................

	// Running server
	app.Run(iris.Addr(":8085"), iris.WithoutServerError(iris.ErrServerClosed))
}

// This the test 2 function
func handler(ctx iris.Context) {
	ctx.WriteString(" Test 2 is working  ")
	ctx.Writef("Hello from method: %s and path: %s", ctx.Method(), ctx.Path())
}

func ApiList(ctx iris.Context) {

	type Medicine struct {
		Id            int    `json:"id"`
		Tradenames    string `json:"Tradenames"`
		Genericnames  string `json:"Genericnames"`
		Chemicalnames string `json:"Chemicalnames"`
		Activationsta string `json:"Activationsta"`
		UsedFor       string `json:"UsedFor"`
		Mrp           string `json:"Mrp"`
	}

	db, dbconnerr := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/medicine")
	rows, err := db.Query("SELECT * FROM medicinelist")
	if dbconnerr != nil {
		fmt.Println("Hmm err :", err.Error())
	}

	var Productlist []Medicine
	for rows.Next() {
		item := Medicine{}
		err2 := rows.Scan(&item.Id, &item.Tradenames, &item.Genericnames, &item.Chemicalnames, &item.Activationsta, &item.UsedFor, &item.Mrp)
		if err2 != nil {
			panic(err2)
		}
		Productlist = append(Productlist, item)
	}
	ctx.JSON(Productlist)
}

func ApiOrder(ctx iris.Context) {
	// ctx.WriteString(" Test 2 is working  ")
	type AutoGenerated struct {
		Id           int    `json:"Id"`
		UserId       int    `json:"UserId"`
		FirstName    string `json:"FirstName"`
		OrderDate    string `json:"OrderDate"`
		DeliveryDate string `json:"DeliveryDate"`
		Products     []struct {
			OrderID    int    `json:"OrderId"`
			Productid  int    `json:"Productid"`
			TradeNames string `json:"TradeNames"`
			Qty        int    `json:"Qty"`
			Price      int    `json:"Price"`
		} `json:"Products"`
	}

	/*
		[{
			"Id": 1,
			"UserId": 2,
			"FirstName": "javed",
			"OrderDate": "12/8/2018",
			"DeliveryDate": "12/9/2018",
			"Products": [{
				"OrderId": 1,
				"Productid": 10,
				"TradeNames": "Napa Extra",
				"Qty": 12,
				"Price": 144
			}, {
				"OrderId": 1,
				"Productid": 10,
				"TradeNames": "Napa Extra",
				"Qty": 12,
				"Price": 144
			}]

		}
		,
		{
			"Id": 1,
			"UserId": 2,
			"FirstName": "javed",
			"OrderDate": "12/8/2018",
			"DeliveryDate": "12/9/2018",
			"Products": [{
				"OrderId": 1,
				"Productid": 10,
				"TradeNames": "Napa Extra",
				"Qty": 12,
				"Price": 144
			}, {
				"OrderId": 1,
				"Productid": 10,
				"TradeNames": "Napa Extra",
				"Qty": 12,
				"Price": 144
			}]

		}
		]

	*/

	// type Product struct {
	// 	OrderId    int
	// 	Productid  int
	// 	TradeNames string
	// 	Qty        int
	// 	Price      int
	// }
	// type Orders struct {
	// 	Id           int
	// 	UserId       int
	// 	FirstName    string
	// 	OrderDate    string
	// 	DeliveryDate string
	// 	Products     []Product
	// }

	// var data []AutoGenerated
	// ctx.ReadJSON(&data)
	// ctx.JSON(data)

	// .............................................................
	// type AutoGenerated struct {
	// 	Name string `json:"name"`
	// }

	var hmm []AutoGenerated

	ctx.ReadJSON(&hmm)

	// fmt.Println("Is it working ")
	// fmt.Println(hmm)
	fmt.Println(hmm)
	ctx.JSON(hmm)
}

// .................................................................. Desclamiers ......................................................

/*
1. Need to Add JWT
2. Need to Add Validation
3. Need to Refactoring The Code

Note : This Ends heare   ...... I will Mode In Another Project Now !!!!

*/
