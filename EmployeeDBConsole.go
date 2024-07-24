package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

type Employee struct {
	ID         int
	Name       string
	Department string
	Salary     float64
	Age        int
	Position   string
}

var db *sql.DB

func main() {
	run()
}

func run() {
	for {
		fmt.Println("\nMenu:")
		fmt.Println("1. Connect to Database")
		fmt.Println("2. Create Table")
		fmt.Println("3. Insert Default Data")
		fmt.Println("4. Get All Employees")
		fmt.Println("5. Find Employee by ID")
		fmt.Println("6. Drop Table")
		fmt.Println("7. Exit")
		fmt.Print("Choose an option: ")

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			connectToDB()
		case 2:
			if db != nil {
				createTable()
			} else {
				fmt.Println("Please connect to the database first.")
			}
		case 3:
			if db != nil {
				insertEmployees()
			} else {
				fmt.Println("Please connect to the database first.")
			}
		case 4:
			if db != nil {
				getAllEmployees()
			} else {
				fmt.Println("Please connect to the database first.")
			}
		case 5:
			if db != nil {
				getEmployeeByID()
			} else {
				fmt.Println("Please connect to the database first.")
			}
		case 6:
			if db != nil {
				dropTable()
			} else {
				fmt.Println("Please connect to the database first.")
			}
		case 7:
			if db != nil {
				closeDBConn()
			}
			fmt.Println("Exiting the program.")
			os.Exit(0)
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}

func connectToDB() {
	var err error
	db, err = ConnectToDB()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	fmt.Println("Database connection established.")
}

func closeDBConn() {
	if db != nil {
		if err := db.Close(); err != nil {
			log.Fatalf("Error closing database connection: %v", err)
		}
		fmt.Println("Database connection closed.")
	}
}

func createTable() {
	if err := CreateTableIfNotExists(db); err != nil {
		log.Fatalf("Error creating table: %v", err)
	}
	fmt.Println("Table created.")
}

func insertEmployees() {
	if err := InsertEmployees(db); err != nil {
		log.Fatalf("Error inserting employees: %v", err)
	}
	fmt.Println("Data inserted into the table.")
}

func getAllEmployees() {
	employees, err := GetAllEmployees(db)
	if err != nil {
		log.Fatalf("Error retrieving all employees: %v", err)
	}
	PrintEmployees(employees)
}

func getEmployeeByID() {
	var id int
	fmt.Print("Enter Employee ID: ")
	fmt.Scanln(&id)

	employee, err := GetEmployeeByID(db, id)
	if err != nil {
		log.Fatalf("Error retrieving employee by ID: %v", err)
	}
	PrintEmployee(employee)
}

func dropTable() {
	if err := DropTable(db); err != nil {
		log.Fatalf("Error dropping table: %v", err)
	}
	fmt.Println("Table dropped.")
}

func ConnectToDB() (*sql.DB, error) {
	connStr := "user=postgres password=hy.par2004 dbname=my_db sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	fmt.Println("Connected to database")
	return db, nil
}

func CreateTableIfNotExists(db *sql.DB) error {
	query := `
    CREATE TABLE IF NOT EXISTS employees (
        id SERIAL PRIMARY KEY,
        name VARCHAR(50) NOT NULL,
        department VARCHAR(50) NOT NULL,
        salary NUMERIC(10, 2) NOT NULL,
        age INTEGER NOT NULL,
        position VARCHAR(50) NOT NULL
    );`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func InsertEmployees(db *sql.DB) error {
	query := `
    INSERT INTO employees (name, department, salary, age, position)
    VALUES
        ('John Doe', 'Marketing', 50000.00, 35, 'Back-End'),
        ('Jane Smith', 'Marketing', 55000.00, 42, 'Designer'),
        ('Bob Johnson', 'IT', 60000.00, 28, 'Back-End'),
        ('Sara Lee', 'IT', 65000.00, 52, 'Front-End'),
        ('Mike Williams', 'HR', 45000.00, 31, 'Q&A'),
        ('Emily Davis', 'HR', 48000.00, 27, 'Q&A'),
        ('David Brown', 'Finance', 70000.00, 47, 'Back-End'),
        ('Samantha Wilson', 'Finance', 75000.00, 55, 'Back-End'),
        ('Tom Garcia', 'Marketing', 52000.00, 29, 'Back-End'),
        ('Olivia Hernandez', 'IT', 62000.00, 38, 'Front-End');
    `
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func GetAllEmployees(db *sql.DB) ([]Employee, error) {
	rows, err := db.Query("SELECT id, name, department, salary, age, position FROM employees;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []Employee
	for rows.Next() {
		var e Employee
		err = rows.Scan(&e.ID, &e.Name, &e.Department, &e.Salary, &e.Age, &e.Position)
		if err != nil {
			fmt.Println(err)
			continue
		}
		employees = append(employees, e)
	}
	return employees, nil
}

func GetEmployeeByID(db *sql.DB, id int) (Employee, error) {
	var e Employee
	row := db.QueryRow("SELECT id, name, department, salary, age, position FROM employees WHERE id = $1", id)
	err := row.Scan(&e.ID, &e.Name, &e.Department, &e.Salary, &e.Age, &e.Position)
	if err != nil {
		return Employee{}, err
	}
	return e, nil
}

func DropTable(db *sql.DB) error {
	query := "DROP TABLE IF EXISTS employees;"
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func PrintEmployees(employees []Employee) {
	fmt.Println("Employees:")
	for _, e := range employees {
		fmt.Printf("%+v\n", e)
	}
}

func PrintEmployee(employee Employee) {
	fmt.Printf("Employee with ID %d: %+v\n", employee.ID, employee)
}
