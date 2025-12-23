package domain

// Character はプレイヤーキャラクターを表す
type Character struct {
	PlayerID   string
	Name       string
	Role       string
	SecretInfo string
	PublicInfo string
}
