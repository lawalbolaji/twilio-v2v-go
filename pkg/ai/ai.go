package ai

type LLM interface {
	GetCompletion(string) (string, error)
}
