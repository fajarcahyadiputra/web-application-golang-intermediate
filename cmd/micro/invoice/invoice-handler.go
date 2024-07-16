package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/signintech/gopdf"
)

// type Order is the type for order
type Order struct {
	ID        int       `json:"id"`
	Quantity  int       `json:"quantity"`
	Amount    int       `json:"amount"`
	Product   string    `json:"product"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func (app *application) CreateAndSendInvoice(w http.ResponseWriter, r *http.Request) {
	// receice json
	var order Order
	// err := app.readJSON(w, r, &order)
	// if err != nil {
	// 	app.badRequest(w, r, err)
	// 	return
	// }
	// generate a pdf invoice
	order.ID = 100
	order.Email = "fajar@cc.coo"
	order.FirstName = "fajar"
	order.LastName = "cp"
	order.Amount = 1000
	order.Quantity = 1
	order.Product = "Widget"
	order.CreatedAt = time.Now()
	err := app.createInvoicePDF(order)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	//send response
	var resp struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}

	resp.Error = false
	resp.Message = fmt.Sprintf("Invoice %d.pdf created and sent to %s", order.ID, order.Email)
	app.writeJSON(w, http.StatusOK, resp)
}

func (app *application) createInvoicePDF(order Order) error {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	// pdf.SetMargins(10, 5, 10, 5)
	pdf.AddPage()
	err := pdf.AddTTFFont("graphik-bold", "./assets/Graphik-Semibold.ttf")
	if err != nil {
		log.Fatal(err)
		return err
	}
	err = pdf.SetFont("graphik-bold", "", 11)
	if err != nil {
		log.Print(err.Error())
		return err
	}
	// t := pdf.ImportPage("./pdf-templates/invoice.pdf", 1, "/MediaBox")
	// // Draw pdf onto page
	// pdf.UseImportedTemplate(t, 0, 0, 215.9, 0)

	//write info
	pdf.SetY(50)
	pdf.SetX(10)
	pdf.SetFont("Times", "", 11)
	pdf.CellWithOption(&gopdf.Rect{W: 97, H: 8}, fmt.Sprintf("Attention: %s %s", order.FirstName, order.LastName), gopdf.CellOption{
		Border:      gopdf.ContentTypeText,
		Align:       gopdf.Left,
		BreakOption: &gopdf.DefaultBreakOption,
	})
	pdf.Br(5)
	pdf.CellWithOption(&gopdf.Rect{W: 97, H: 8}, order.Email, gopdf.CellOption{
		Border:      gopdf.ContentTypeText,
		Align:       gopdf.Left,
		BreakOption: &gopdf.DefaultBreakOption,
	})
	pdf.Br(5)
	pdf.CellWithOption(&gopdf.Rect{W: 97, H: 8}, order.CreatedAt.Format("2006-01-02"), gopdf.CellOption{
		Border:      gopdf.ContentTypeText,
		Align:       gopdf.Left,
		BreakOption: &gopdf.DefaultBreakOption,
	})

	pdf.SetX(58)
	pdf.SetY(93)
	pdf.CellWithOption(&gopdf.Rect{W: 155, H: 8}, order.Product, gopdf.CellOption{
		Border:      gopdf.ContentTypeText,
		Align:       gopdf.Left,
		BreakOption: &gopdf.DefaultBreakOption,
	})

	pdf.SetX(166)
	pdf.CellWithOption(&gopdf.Rect{W: 20, H: 8}, fmt.Sprintf("%d", order.Quantity), gopdf.CellOption{
		Border:      gopdf.ContentTypeText,
		Align:       gopdf.Left,
		BreakOption: &gopdf.DefaultBreakOption,
	})

	pdf.SetX(185)
	pdf.CellWithOption(&gopdf.Rect{W: 20, H: 8}, fmt.Sprintf("$%.2f", float32(order.Quantity/100.0)), gopdf.CellOption{
		Border:      gopdf.ContentTypeText,
		Align:       gopdf.Right,
		BreakOption: &gopdf.DefaultBreakOption,
	})

	invoicePath := fmt.Sprintf("./invoices/%d.pdf", order.ID)
	err = pdf.WritePdf(invoicePath)
	if err != nil {
		return err
	}

	return nil
}
