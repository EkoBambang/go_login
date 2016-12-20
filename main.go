package main

import "net/http"
import "database/sql"
import _ "github.com/go-sql-driver/mysql"
import "github.com/gorilla/mux"
import "github.com/gorilla/securecookie"
import "html/template"
import "os"
import "fmt"

func connect() *sql.DB {
	var db, err = sql.Open("mysql", "root:@/login_go")
	err = db.Ping()
	if err != nil {
		os.Exit(0)
	}
	return db
}

//---------------------------html dalam----------------------------------------//
var router = mux.NewRouter()

var cookiehandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

const login = `<html>
            <form method="POST" action="/mau_login">
			username :<br><input type="text" name="username"><br>
			password :<br><input type="password" name="password"><br>
			<input type="submit" name="login">
			</form?
               </html>`
const index = `<html>
				<p>selamat datang {{.nama}},anda baru saja login</p><br>
				<form method="post" action="/mau_logout">
				
				<button type="submit">Logout</button>
				</form>
				</html>`
const gagal = `<html>
               data yang anda masukan salah<br>
			   <a href="/"> kembali</a>
               </html>`

//---------------------menu login----------------------------------------------//
func login_menu(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, login)
}

func halaman_index(res http.ResponseWriter, req *http.Request) {
	akses := namauser(req)
	if akses != "" {

		halaman, _ := template.New("halaman").Parse(index)

		isi := map[string]string{
			"nama": akses,
		}
		halaman.Execute(res, isi)

	} else {
		http.Redirect(res, req, "/", 301)
	}
}

func gagal_login(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, gagal)
}

func mau_logout(res http.ResponseWriter, req *http.Request) {
	clear_session(res)
	http.Redirect(res, req, "/", 301)
}

func mau_login(res http.ResponseWriter, req *http.Request) {
	db := connect()
	defer db.Close()

	username_login := req.FormValue("username")
	pass_login := req.FormValue("password")
	jalur := "/"

	var nama, username, pass string

	db.QueryRow("select * from user where username=?", username_login).Scan(&nama, &username, &pass)

	if pass == pass_login {
		setsesi(nama, res)
		jalur = "/index"
	} else {
		jalur = "/gagal"
	}
	http.Redirect(res, req, jalur, 302)

}

//---------------------------sesi----------------------------------------------//

func setsesi(nama_user string, res http.ResponseWriter) {
	value := map[string]string{
		"name": nama_user,
	}
	if encoded, err := cookiehandler.Encode("session", value); err == nil {
		cookie_ku := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(res, cookie_ku)
	}
}

func namauser(req *http.Request) (nama_usernya string) {
	if cookie_ini, err := req.Cookie("session"); err == nil {
		nilai_cookie := make(map[string]string)
		if err = cookiehandler.Decode("session", cookie_ini.Value, &nilai_cookie); err == nil {
			nama_usernya = nilai_cookie["name"]
		}

	}
	return nama_usernya
}

func clear_session(res http.ResponseWriter) {
	bersihkan_cookie_ku := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(res, bersihkan_cookie_ku)
}

//----------------------------------------main function-----------------------//
func main() {
	router.HandleFunc("/", login_menu)
	router.HandleFunc("/mau_login", mau_login)
	router.HandleFunc("/index", halaman_index)
	router.HandleFunc("/gagal", gagal_login)
	router.HandleFunc("/mau_logout", mau_logout)

	http.Handle("/", router)

	fmt.Println("menjalankan server via localhost...")
	http.ListenAndServe(":8080", nil)
}
