package storage

import (
    "fmt"
    "library-backend/internal/models"
    "time"
)

func CreateBook(b models.Book) error {
    _, err := DB.Exec(`
        INSERT INTO books (id, title, author, isbn, year, status) 
        VALUES (?, ?, ?, ?, ?, ?)`,
        b.ID, b.Title, b.Author, b.ISBN, b.Year, b.Status)
    return err
}

func GetBookByID(id string) (models.Book, error) {
    var b models.Book
    query := `SELECT id, title, author, isbn, year, status FROM books WHERE id = ?`
    err := DB.QueryRow(query, id).Scan(&b.ID, &b.Title, &b.Author, &b.ISBN, &b.Year, &b.Status)
    return b, err
}

func UpdateBook(b models.Book) error {
    _, err := DB.Exec(`
        UPDATE books 
        SET title = ?, author = ?, isbn = ?, year = ?, status = ? 
        WHERE id = ?`,
        b.Title, b.Author, b.ISBN, b.Year, b.Status, b.ID)
    return err
}

func ListBooks(limit, offset int, author, status string) ([]models.Book, error) {
    query := `SELECT id, title, author, isbn, year, status FROM books WHERE 1=1`
    args := []interface{}{}

    if author != "" {
        query += " AND author = ?"
        args = append(args, author)
    }
    if status != "" {
        query += " AND status = ?"
        args = append(args, status)
    }

    query += " LIMIT ? OFFSET ?"
    args = append(args, limit, offset)

    rows, err := DB.Query(query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var books []models.Book
    for rows.Next() {
        var b models.Book
        err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.ISBN, &b.Year, &b.Status)
        if err != nil {
            return nil, err
        }
        books = append(books, b)
    }
    return books, nil
}

func UpdateBookStatus(bookID, status string) error {
    _, err := DB.Exec(`UPDATE books SET status = ? WHERE id = ?`, status, bookID)
    return err
}

func CreateUser(u models.User) error {
    _, err := DB.Exec(`
        INSERT INTO users (id, name, email, password, role, registration_date) 
        VALUES (?, ?, ?, ?, ?, ?)`,
        u.ID, u.Name, u.Email, u.Password, u.Role, u.RegistrationDate)
    return err
}

func GetUserByEmail(email string) (models.User, error) {
    var u models.User
    query := `SELECT id, name, email, password, registration_date, role FROM users WHERE email = ?`
    err := DB.QueryRow(query, email).Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.RegistrationDate, &u.Role)
    return u, err
}

func GetUserByID(id string) (models.User, error) {
    var u models.User
    query := `SELECT id, name, email, registration_date, role FROM users WHERE id = ?`
    err := DB.QueryRow(query, id).Scan(&u.ID, &u.Name, &u.Email, &u.RegistrationDate, &u.Role)
    return u, err
}

func GetUserBooks(userID string) ([]models.Book, error) {
    query := `
        SELECT b.id, b.title, b.author, b.isbn, b.year, b.status
        FROM books b
        JOIN issues i ON b.id = i.book_id
        WHERE i.user_id = ? AND i.return_date IS NULL`
    
    rows, err := DB.Query(query, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var books []models.Book
    for rows.Next() {
        var b models.Book
        err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.ISBN, &b.Year, &b.Status)
        if err != nil {
            return nil, err
        }
        books = append(books, b)
    }
    return books, nil
}

func CreateIssue(issue models.Issue) error {
    _, err := DB.Exec(`
        INSERT INTO issues (book_id, user_id, due_date) 
        VALUES (?, ?, ?)`,
        issue.BookID, issue.UserID, issue.DueDate)
    return err
}

func GetActiveIssueByBookID(bookID string) (models.Issue, error) {
    var i models.Issue
    query := `SELECT id, book_id, user_id, issue_date, due_date, return_date 
              FROM issues 
              WHERE book_id = ? AND return_date IS NULL`
    err := DB.QueryRow(query, bookID).Scan(&i.ID, &i.BookID, &i.UserID, &i.IssueDate, &i.DueDate, &i.ReturnDate)
    return i, err
}

func GetIssueByID(issueID int) (models.Issue, error) {
    var i models.Issue
    query := `SELECT id, book_id, user_id, issue_date, due_date, return_date 
              FROM issues WHERE id = ?`
    err := DB.QueryRow(query, issueID).Scan(&i.ID, &i.BookID, &i.UserID, &i.IssueDate, &i.DueDate, &i.ReturnDate)
    return i, err
}

func ReturnBook(issueID int, returnDate time.Time) error {
    _, err := DB.Exec(`UPDATE issues SET return_date = ? WHERE id = ?`, returnDate, issueID)
    return err
}

func CheckBookAvailability(bookID string) (bool, error) {
    var status string
    query := `SELECT status FROM books WHERE id = ?`
    err := DB.QueryRow(query, bookID).Scan(&status)
    if err != nil {
        return false, err
    }
    return status == "Available", nil
}

func ValidateUserExists(userID string) error {
    var exists bool
    query := `SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)`
    err := DB.QueryRow(query, userID).Scan(&exists)
    if err != nil {
        return err
    }
    if !exists {
        return fmt.Errorf("user not found")
    }
    return nil
}
