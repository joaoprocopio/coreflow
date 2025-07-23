package propostas

type PropostaResponse struct {
	ID          int32                 `json:"id"`
	Status      PropostaStatus        `json:"status"`
	Name        string                `json:"name"`
	Assignee    *User                 `json:"assignee"`
	Attachments []*PropostaAttachment `json:"attachments"`
}

func SerializePropostasToResponse(propostas []*Proposta) []PropostaResponse {
	response := make([]PropostaResponse, len(propostas))

	for i, proposta := range propostas {
		attachments := proposta.Attachments
		if attachments == nil {
			attachments = []*PropostaAttachment{}
		}

		response[i] = PropostaResponse{
			ID:          proposta.ID,
			Status:      proposta.Status,
			Name:        proposta.Name,
			Assignee:    proposta.Assignee,
			Attachments: attachments,
		}
	}

	return response
}
