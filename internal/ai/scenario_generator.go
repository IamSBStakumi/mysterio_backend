package ai

import (
	"context"
	"fmt"

	"github.com/IamSBStakumi/mysterio_backend/internal/domain"
)

// ScenarioGenerator はシナリオを生成するインターフェース
type ScenarioGenerator interface {
	Generate(ctx context.Context, playerCount int, difficulty domain.Difficulty) (*domain.Scenario, error)
}

// MockScenarioGenerator はモックのシナリオジェネレーター（開発用）
type MockScenarioGenerator struct{}

// NewMockScenarioGenerator は新しいモックジェネレーターを作成する
func NewMockScenarioGenerator() *MockScenarioGenerator {
	return &MockScenarioGenerator{}
}

// Generate はモックシナリオを生成する
func (g *MockScenarioGenerator) Generate(
	ctx context.Context,
	playerCount int,
	difficulty domain.Difficulty,
) (*domain.Scenario, error) {
	duration := g.getDuration(difficulty)
	phaseCount := duration / 20

	// プレイヤー数に応じたキャラクターを生成
	characters := g.generateCharacters(playerCount)

	// フェーズを生成
	phases := g.generatePhases(phaseCount)

	scenario := &domain.Scenario{
		Background: "豪華客船「オーシャンスター号」の船上パーティーで、有名な宝石商が殺害された。" +
			"乗客の中に犯人がいる。船は嵐で航行不能になり、警察の到着まで6時間。" +
			"あなたたちは真相を暴かなければならない。",
		Setting: "豪華客船のVIPラウンジ。外は嵐。時刻は深夜2時。",
		Characters: characters,
		Truth: fmt.Sprintf("犯人は%sである。動機は過去の取引での裏切り。" +
			"凶器はラウンジにあったペーパーナイフ。被害者を船室に呼び出し、" +
			"口論の末に刺殺した。", characters[0].Name),
		Phases:     phases,
		Difficulty: difficulty,
	}

	return scenario, nil
}

func (g *MockScenarioGenerator) generateCharacters(count int) []domain.Character {
	characterTemplates := []struct {
		name       string
		role       string
		secretInfo string
		publicInfo string
	}{
		{
			name: "アレックス・クロフォード",
			role: "宝石商のビジネスパートナー",
			secretInfo: "被害者との間に多額の借金がある。返済を迫られていた。" +
				"事件当夜、被害者の船室に向かうところを目撃されないよう、" +
				"別ルートを使った。実は犯人である。",
			publicInfo: "被害者とは20年来のビジネスパートナー。表向きは良好な関係。",
		},
		{
			name: "エミリー・ハートウェル",
			role: "被害者の秘書",
			secretInfo: "被害者から不当な扱いを受けていた。解雇を言い渡される直前だった。" +
				"事件当夜、ラウンジで被害者と口論していた。",
			publicInfo: "5年間、被害者の秘書として働いている。真面目で几帳面な性格。",
		},
		{
			name: "ダニエル・ヴァンクリーフ",
			role: "ライバル宝石商",
			secretInfo: "被害者に大きな取引を横取りされ、会社が倒産寸前。" +
				"復讐の機会を狙っていた。パーティーには招待されていないが、" +
				"偽名で乗船していた。",
			publicInfo: "同業者だが、被害者とは犬猿の仲として知られている。",
		},
		{
			name: "ソフィア・ブラックウッド",
			role: "被害者の元妻",
			secretInfo: "離婚後も財産分与で争っていた。被害者が隠し財産を持っていると" +
				"睨んでいる。事件当夜、被害者を尾行していた。",
			publicInfo: "3年前に離婚。子供の親権を持っている。パーティーには招待客として参加。",
		},
		{
			name: "マーカス・レインズ",
			role: "船の警備主任",
			secretInfo: "過去に被害者から賄賂を受け取り、密輸を見逃したことがある。" +
				"それが発覚することを恐れていた。被害者から脅迫されていた。",
			publicInfo: "10年のキャリアを持つベテラン警備員。規則に厳格と評判。",
		},
	}

	characters := make([]domain.Character, 0, count)
	for i := 0; i < count && i < len(characterTemplates); i++ {
		template := characterTemplates[i]
		characters = append(characters, domain.Character{
			PlayerID:   "", // サービス層で設定される
			Name:       template.name,
			Role:       template.role,
			SecretInfo: template.secretInfo,
			PublicInfo: template.publicInfo,
		})
	}

	return characters
}

func (g *MockScenarioGenerator) generatePhases(count int) []domain.Phase {
	basePhases := []domain.Phase{
		{
			Number:      0,
			Type:        domain.PhaseTypeIntroduction,
			Description: "事件の概要と各自の立場を確認",
			PublicText: "被害者：サミュエル・ゴールドスタイン（65歳）、有名宝石商。" +
				"死亡推定時刻：深夜1時30分。死因：鋭利な刃物による刺殺。" +
				"発見場所：被害者の船室。",
			Duration: 10,
		},
		{
			Number:      1,
			Type:        domain.PhaseTypeDiscussion,
			Description: "第一回討議 - 各自のアリバイを確認",
			PublicText: "事件当時、あなたはどこで何をしていましたか？" +
				"互いのアリバイを確認しましょう。",
			Duration: 15,
		},
		{
			Number:      2,
			Type:        domain.PhaseInvestigation1,
			Description: "証拠品の検証",
			PublicText: "発見された証拠品：\n" +
				"1. 血痕のついたペーパーナイフ（凶器）\n" +
				"2. 被害者の手帳（最後のページに「A.C. 2:00」と書かれている）\n" +
				"3. 廊下の監視カメラ映像（一部が消去されている）",
			Duration: 20,
		},
		{
			Number:      3,
			Type:        domain.PhaseTypeDiscussion,
			Description: "第二回討議 - 動機と証拠の分析",
			PublicText: "証拠品を踏まえて、誰に最も動機があるか議論しましょう。" +
				"矛盾点や不審な点を指摘してください。",
			Duration: 20,
		},
		{
			Number:      4,
			Type:        domain.PhaseTypeVoting,
			Description: "最終投票 - 犯人の決定",
			PublicText: "これまでの議論を踏まえ、犯人だと思う人物に投票してください。" +
				"最も票を集めた人物が犯人として告発されます。",
			Duration: 10,
		},
		{
			Number:      5,
			Type:        domain.PhaseTypeReveal,
			Description: "真相の公開",
			PublicText: "真相が明かされます。",
			Duration: 5,
		},
	}

	// フェーズ数に応じて調整（最小3、最大6）
	if count < 3 {
		count = 3
	}
	if count > 6 {
		count = 6
	}

	phases := make([]domain.Phase, 0, count)
	for i := 0; i < count && i < len(basePhases); i++ {
		phases = append(phases, basePhases[i])
	}

	return phases
}

func (g *MockScenarioGenerator) getDuration(difficulty domain.Difficulty) int {
	switch difficulty {
	case domain.DifficultyEasy:
		return 60
	case domain.DifficultyMedium:
		return 90
	case domain.DifficultyHard:
		return 120
	default:
		return 90
	}
}

// AnthropicScenarioGenerator は実際のAnthropic APIを使用するジェネレーター
// TODO: Anthropic APIが利用可能になったら実装
type AnthropicScenarioGenerator struct {
	apiKey string
}

// NewAnthropicScenarioGenerator は新しいAnthropicジェネレーターを作成する
func NewAnthropicScenarioGenerator(apiKey string) *AnthropicScenarioGenerator {
	return &AnthropicScenarioGenerator{
		apiKey: apiKey,
	}
}

// Generate はAnthropic APIを使用してシナリオを生成する
func (g *AnthropicScenarioGenerator) Generate(
	ctx context.Context,
	playerCount int,
	difficulty domain.Difficulty,
) (*domain.Scenario, error) {
	// TODO: Anthropic API呼び出しを実装
	return nil, fmt.Errorf("not implemented yet")
}
