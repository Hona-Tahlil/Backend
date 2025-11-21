package enums

type RequestStatus uint

const (
	Pending RequestStatus = iota + 1
	Accepted
	Paid
	Finished
	Canceled
	Dismissed
	Conflict
	// paid and canceled or unpaid
	// user canceled or pet sitter
)

func (requestStatus RequestStatus) String() string {
	switch requestStatus {
	case RequestStatus(Pending):
		return "درخواست در انتظار تایید است"
	case RequestStatus(Accepted):
		return "درخواست تایید شده و در انتظار پرداخت است"
	case RequestStatus(Paid):
		return "درخواست پرداخت و نهایی شده"
	case RequestStatus(Finished):
		return "خدمت مورد نظر تمام شده است"
	case RequestStatus(Canceled):
		return "درخواست کنسل شده است"
	case RequestStatus(Dismissed):
		return "درخواست توسط پت سیتر رد شده است"
	case RequestStatus(Conflict):
		return "این درخواست با درخواست دیگری از پتیار همپوشانی دارد"
	}
	return ""
}
