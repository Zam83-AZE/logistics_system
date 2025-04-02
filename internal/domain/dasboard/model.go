package dashboard

// Summary dashboard üçün əsas statistika məlumatlarını təmsil edir
type Summary struct {
	TotalCustomers  int `db:"total_customers"`
	TotalContainers int `db:"total_containers"`
	ActiveShipments int `db:"active_shipments"`
	PendingInvoices int `db:"pending_invoices"`
}

// DashboardData dashboard üçün bütün lazımi məlumatları təmsil edir
type DashboardData struct {
	Summary     Summary
	UserName    string
	CurrentPage string
	Error       string
}
