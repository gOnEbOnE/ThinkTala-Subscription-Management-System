package models

import "time"

type SupportTicket struct {
	ID            string    `json:"id"`
	UserName      string    `json:"user_name"`
	UserEmail     string    `json:"user_email,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	Title         string    `json:"title"`
	Status        string    `json:"status"`
	Category      string    `json:"category,omitempty"`
	Description   string    `json:"description,omitempty"`
	HasAttachment bool      `json:"has_attachment"`
	AttachmentURL string    `json:"attachment_url,omitempty"`
}

type CreateSupportTicketRequest struct {
	Title       string `json:"title"`
	Category    string `json:"category"`
	Description string `json:"description"`
}

type AuthenticatedUser struct {
	ID    string
	Name  string
	Email string
	Role  string
}

type SupportTicketPagination struct {
	Page       int  `json:"page"`
	PerPage    int  `json:"per_page"`
	Total      int  `json:"total"`
	TotalPages int  `json:"total_pages"`
	HasNext    bool `json:"has_next"`
	HasPrev    bool `json:"has_prev"`
}

type SupportTicketsResponse struct {
	Data       []SupportTicket         `json:"data"`
	Pagination SupportTicketPagination `json:"pagination"`
}

type SupportTicketReply struct {
	ID        string    `json:"id"`
	TicketID  string    `json:"ticket_id"`
	AdminID   string    `json:"admin_id"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

type SupportTicketDetailResponse struct {
	Ticket  SupportTicket        `json:"ticket"`
	Replies []SupportTicketReply `json:"replies"`
}

type CreateSupportTicketReplyRequest struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

type CreateSupportTicketInput struct {
	Title          string
	Category       string
	Description    string
	AttachmentName string
	AttachmentMIME string
	AttachmentData []byte
}
