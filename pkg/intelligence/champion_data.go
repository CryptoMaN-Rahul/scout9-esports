package intelligence

// ChampionData contains metadata about a League of Legends champion
// Data sourced from official League of Legends Wiki (wiki.leagueoflegends.com)
type ChampionData struct {
	Name       string   `json:"name"`
	Class      string   `json:"class"`      // Primary class: Controller, Fighter, Mage, Marksman, Slayer, Tank, Specialist
	Subclass   string   `json:"subclass"`   // Subclass: Enchanter, Catcher, Juggernaut, Diver, Burst, Battlemage, Artillery, Assassin, Skirmisher, Vanguard, Warden
	Positions  []string `json:"positions"`  // Common positions: Top, Jungle, Mid, Bot, Support
	DamageType string   `json:"damageType"` // Physical, Magic, Mixed
	Range      string   `json:"range"`      // Melee, Ranged
}

// ChampionDatabase is the comprehensive database of all LoL champions with their classes
// Based on official Riot Games classification from wiki.leagueoflegends.com/en-us/Champion_classes
var ChampionDatabase = map[string]ChampionData{
	// Controllers - Enchanters
	"Janna":    {Name: "Janna", Class: "Controller", Subclass: "Enchanter", Positions: []string{"Support"}, DamageType: "Magic", Range: "Ranged"},
	"Karma":    {Name: "Karma", Class: "Controller", Subclass: "Enchanter", Positions: []string{"Support", "Mid"}, DamageType: "Magic", Range: "Ranged"},
	"Lulu":     {Name: "Lulu", Class: "Controller", Subclass: "Enchanter", Positions: []string{"Support"}, DamageType: "Magic", Range: "Ranged"},
	"Milio":    {Name: "Milio", Class: "Controller", Subclass: "Enchanter", Positions: []string{"Support"}, DamageType: "Magic", Range: "Ranged"},
	"Nami":     {Name: "Nami", Class: "Controller", Subclass: "Enchanter", Positions: []string{"Support"}, DamageType: "Magic", Range: "Ranged"},
	"Renata":   {Name: "Renata Glasc", Class: "Controller", Subclass: "Enchanter", Positions: []string{"Support"}, DamageType: "Magic", Range: "Ranged"},
	"Seraphine": {Name: "Seraphine", Class: "Controller", Subclass: "Enchanter", Positions: []string{"Support", "Mid", "Bot"}, DamageType: "Magic", Range: "Ranged"},
	"Sona":     {Name: "Sona", Class: "Controller", Subclass: "Enchanter", Positions: []string{"Support"}, DamageType: "Magic", Range: "Ranged"},
	"Soraka":   {Name: "Soraka", Class: "Controller", Subclass: "Enchanter", Positions: []string{"Support"}, DamageType: "Magic", Range: "Ranged"},
	"Yuumi":    {Name: "Yuumi", Class: "Controller", Subclass: "Enchanter", Positions: []string{"Support"}, DamageType: "Magic", Range: "Ranged"},
	"Ivern":    {Name: "Ivern", Class: "Controller", Subclass: "Enchanter", Positions: []string{"Jungle"}, DamageType: "Magic", Range: "Ranged"},
	"Taric":    {Name: "Taric", Class: "Controller", Subclass: "Enchanter", Positions: []string{"Support"}, DamageType: "Magic", Range: "Melee"},

	// Controllers - Catchers
	"Bard":       {Name: "Bard", Class: "Controller", Subclass: "Catcher", Positions: []string{"Support"}, DamageType: "Magic", Range: "Ranged"},
	"Blitzcrank": {Name: "Blitzcrank", Class: "Controller", Subclass: "Catcher", Positions: []string{"Support"}, DamageType: "Magic", Range: "Melee"},
	"Morgana":    {Name: "Morgana", Class: "Controller", Subclass: "Catcher", Positions: []string{"Support", "Mid"}, DamageType: "Magic", Range: "Ranged"},
	"Neeko":      {Name: "Neeko", Class: "Controller", Subclass: "Catcher", Positions: []string{"Mid", "Support"}, DamageType: "Magic", Range: "Ranged"},
	"Pyke":       {Name: "Pyke", Class: "Controller", Subclass: "Catcher", Positions: []string{"Support"}, DamageType: "Physical", Range: "Melee"},
	"Rakan":      {Name: "Rakan", Class: "Controller", Subclass: "Catcher", Positions: []string{"Support"}, DamageType: "Magic", Range: "Melee"},
	"Thresh":     {Name: "Thresh", Class: "Controller", Subclass: "Catcher", Positions: []string{"Support"}, DamageType: "Magic", Range: "Melee"},
	"Zyra":       {Name: "Zyra", Class: "Controller", Subclass: "Catcher", Positions: []string{"Support"}, DamageType: "Magic", Range: "Ranged"},
	"Nautilus":   {Name: "Nautilus", Class: "Controller", Subclass: "Catcher", Positions: []string{"Support"}, DamageType: "Magic", Range: "Melee"},

	// Fighters - Juggernauts
	"Aatrox":      {Name: "Aatrox", Class: "Fighter", Subclass: "Juggernaut", Positions: []string{"Top"}, DamageType: "Physical", Range: "Melee"},
	"Darius":      {Name: "Darius", Class: "Fighter", Subclass: "Juggernaut", Positions: []string{"Top"}, DamageType: "Physical", Range: "Melee"},
	"Dr. Mundo":   {Name: "Dr. Mundo", Class: "Fighter", Subclass: "Juggernaut", Positions: []string{"Top", "Jungle"}, DamageType: "Magic", Range: "Melee"},
	"Garen":       {Name: "Garen", Class: "Fighter", Subclass: "Juggernaut", Positions: []string{"Top"}, DamageType: "Physical", Range: "Melee"},
	"Illaoi":      {Name: "Illaoi", Class: "Fighter", Subclass: "Juggernaut", Positions: []string{"Top"}, DamageType: "Physical", Range: "Melee"},
	"Mordekaiser": {Name: "Mordekaiser", Class: "Fighter", Subclass: "Juggernaut", Positions: []string{"Top"}, DamageType: "Magic", Range: "Melee"},
	"Nasus":       {Name: "Nasus", Class: "Fighter", Subclass: "Juggernaut", Positions: []string{"Top"}, DamageType: "Physical", Range: "Melee"},
	"Sett":        {Name: "Sett", Class: "Fighter", Subclass: "Juggernaut", Positions: []string{"Top", "Support"}, DamageType: "Physical", Range: "Melee"},
	"Shyvana":     {Name: "Shyvana", Class: "Fighter", Subclass: "Juggernaut", Positions: []string{"Jungle"}, DamageType: "Mixed", Range: "Melee"},
	"Sion":        {Name: "Sion", Class: "Fighter", Subclass: "Juggernaut", Positions: []string{"Top"}, DamageType: "Physical", Range: "Melee"},
	"Trundle":     {Name: "Trundle", Class: "Fighter", Subclass: "Juggernaut", Positions: []string{"Top", "Jungle"}, DamageType: "Physical", Range: "Melee"},
	"Udyr":        {Name: "Udyr", Class: "Fighter", Subclass: "Juggernaut", Positions: []string{"Jungle", "Top"}, DamageType: "Mixed", Range: "Melee"},
	"Urgot":       {Name: "Urgot", Class: "Fighter", Subclass: "Juggernaut", Positions: []string{"Top"}, DamageType: "Physical", Range: "Ranged"},
	"Volibear":    {Name: "Volibear", Class: "Fighter", Subclass: "Juggernaut", Positions: []string{"Top", "Jungle"}, DamageType: "Mixed", Range: "Melee"},
	"Yorick":      {Name: "Yorick", Class: "Fighter", Subclass: "Juggernaut", Positions: []string{"Top"}, DamageType: "Physical", Range: "Melee"},

	// Fighters - Divers
	"Camille":    {Name: "Camille", Class: "Fighter", Subclass: "Diver", Positions: []string{"Top"}, DamageType: "Mixed", Range: "Melee"},
	"Diana":      {Name: "Diana", Class: "Fighter", Subclass: "Diver", Positions: []string{"Jungle", "Mid"}, DamageType: "Magic", Range: "Melee"},
	"Elise":      {Name: "Elise", Class: "Fighter", Subclass: "Diver", Positions: []string{"Jungle"}, DamageType: "Magic", Range: "Ranged"},
	"Hecarim":    {Name: "Hecarim", Class: "Fighter", Subclass: "Diver", Positions: []string{"Jungle"}, DamageType: "Physical", Range: "Melee"},
	"Irelia":     {Name: "Irelia", Class: "Fighter", Subclass: "Diver", Positions: []string{"Top", "Mid"}, DamageType: "Physical", Range: "Melee"},
	"Jarvan IV":  {Name: "Jarvan IV", Class: "Fighter", Subclass: "Diver", Positions: []string{"Jungle"}, DamageType: "Physical", Range: "Melee"},
	"Kled":       {Name: "Kled", Class: "Fighter", Subclass: "Diver", Positions: []string{"Top"}, DamageType: "Physical", Range: "Melee"},
	"Lee Sin":    {Name: "Lee Sin", Class: "Fighter", Subclass: "Diver", Positions: []string{"Jungle"}, DamageType: "Physical", Range: "Melee"},
	"Olaf":       {Name: "Olaf", Class: "Fighter", Subclass: "Diver", Positions: []string{"Top", "Jungle"}, DamageType: "Physical", Range: "Melee"},
	"Pantheon":   {Name: "Pantheon", Class: "Fighter", Subclass: "Diver", Positions: []string{"Top", "Mid", "Support"}, DamageType: "Physical", Range: "Melee"},
	"Rek'Sai":    {Name: "Rek'Sai", Class: "Fighter", Subclass: "Diver", Positions: []string{"Jungle"}, DamageType: "Physical", Range: "Melee"},
	"Renekton":   {Name: "Renekton", Class: "Fighter", Subclass: "Diver", Positions: []string{"Top"}, DamageType: "Physical", Range: "Melee"},
	"Vi":         {Name: "Vi", Class: "Fighter", Subclass: "Diver", Positions: []string{"Jungle"}, DamageType: "Physical", Range: "Melee"},
	"Warwick":    {Name: "Warwick", Class: "Fighter", Subclass: "Diver", Positions: []string{"Jungle", "Top"}, DamageType: "Mixed", Range: "Melee"},
	"Wukong":     {Name: "Wukong", Class: "Fighter", Subclass: "Diver", Positions: []string{"Jungle", "Top"}, DamageType: "Physical", Range: "Melee"},
	"Xin Zhao":   {Name: "Xin Zhao", Class: "Fighter", Subclass: "Diver", Positions: []string{"Jungle"}, DamageType: "Physical", Range: "Melee"},

	// Mages - Burst
	"Ahri":    {Name: "Ahri", Class: "Mage", Subclass: "Burst", Positions: []string{"Mid"}, DamageType: "Magic", Range: "Ranged"},
	"Annie":   {Name: "Annie", Class: "Mage", Subclass: "Burst", Positions: []string{"Mid", "Support"}, DamageType: "Magic", Range: "Ranged"},
	"Brand":   {Name: "Brand", Class: "Mage", Subclass: "Burst", Positions: []string{"Support", "Mid"}, DamageType: "Magic", Range: "Ranged"},
	"Hwei":    {Name: "Hwei", Class: "Mage", Subclass: "Burst", Positions: []string{"Mid", "Support"}, DamageType: "Magic", Range: "Ranged"},
	"LeBlanc": {Name: "LeBlanc", Class: "Mage", Subclass: "Burst", Positions: []string{"Mid"}, DamageType: "Magic", Range: "Ranged"},
	"Lissandra": {Name: "Lissandra", Class: "Mage", Subclass: "Burst", Positions: []string{"Mid"}, DamageType: "Magic", Range: "Ranged"},
	"Lux":     {Name: "Lux", Class: "Mage", Subclass: "Burst", Positions: []string{"Support", "Mid"}, DamageType: "Magic", Range: "Ranged"},
	"Orianna": {Name: "Orianna", Class: "Mage", Subclass: "Burst", Positions: []string{"Mid"}, DamageType: "Magic", Range: "Ranged"},
	"Syndra":  {Name: "Syndra", Class: "Mage", Subclass: "Burst", Positions: []string{"Mid"}, DamageType: "Magic", Range: "Ranged"},
	"Twisted Fate": {Name: "Twisted Fate", Class: "Mage", Subclass: "Burst", Positions: []string{"Mid"}, DamageType: "Magic", Range: "Ranged"},
	"Veigar":  {Name: "Veigar", Class: "Mage", Subclass: "Burst", Positions: []string{"Mid", "Support"}, DamageType: "Magic", Range: "Ranged"},
	"Vex":     {Name: "Vex", Class: "Mage", Subclass: "Burst", Positions: []string{"Mid"}, DamageType: "Magic", Range: "Ranged"},
	"Viktor":  {Name: "Viktor", Class: "Mage", Subclass: "Burst", Positions: []string{"Mid"}, DamageType: "Magic", Range: "Ranged"},
	"Zoe":     {Name: "Zoe", Class: "Mage", Subclass: "Burst", Positions: []string{"Mid"}, DamageType: "Magic", Range: "Ranged"},

	// Mages - Battlemage
	"Anivia":    {Name: "Anivia", Class: "Mage", Subclass: "Battlemage", Positions: []string{"Mid"}, DamageType: "Magic", Range: "Ranged"},
	"Aurelion Sol": {Name: "Aurelion Sol", Class: "Mage", Subclass: "Battlemage", Positions: []string{"Mid"}, DamageType: "Magic", Range: "Ranged"},
	"Cassiopeia": {Name: "Cassiopeia", Class: "Mage", Subclass: "Battlemage", Positions: []string{"Mid"}, DamageType: "Magic", Range: "Ranged"},
	"Karthus":   {Name: "Karthus", Class: "Mage", Subclass: "Battlemage", Positions: []string{"Jungle", "Mid"}, DamageType: "Magic", Range: "Ranged"},
	"Malzahar":  {Name: "Malzahar", Class: "Mage", Subclass: "Battlemage", Positions: []string{"Mid"}, DamageType: "Magic", Range: "Ranged"},
	"Rumble":    {Name: "Rumble", Class: "Mage", Subclass: "Battlemage", Positions: []string{"Top", "Mid"}, DamageType: "Magic", Range: "Melee"},
	"Ryze":      {Name: "Ryze", Class: "Mage", Subclass: "Battlemage", Positions: []string{"Mid", "Top"}, DamageType: "Magic", Range: "Ranged"},
	"Swain":     {Name: "Swain", Class: "Mage", Subclass: "Battlemage", Positions: []string{"Support", "Mid", "Bot"}, DamageType: "Magic", Range: "Ranged"},
	"Taliyah":   {Name: "Taliyah", Class: "Mage", Subclass: "Battlemage", Positions: []string{"Jungle", "Mid"}, DamageType: "Magic", Range: "Ranged"},
	"Vladimir":  {Name: "Vladimir", Class: "Mage", Subclass: "Battlemage", Positions: []string{"Mid", "Top"}, DamageType: "Magic", Range: "Ranged"},

	// Mages - Artillery
	"Jayce":   {Name: "Jayce", Class: "Mage", Subclass: "Artillery", Positions: []string{"Top", "Mid"}, DamageType: "Physical", Range: "Ranged"},
	"Kog'Maw": {Name: "Kog'Maw", Class: "Mage", Subclass: "Artillery", Positions: []string{"Bot"}, DamageType: "Mixed", Range: "Ranged"},
	"Vel'Koz": {Name: "Vel'Koz", Class: "Mage", Subclass: "Artillery", Positions: []string{"Support", "Mid"}, DamageType: "Magic", Range: "Ranged"},
	"Xerath":  {Name: "Xerath", Class: "Mage", Subclass: "Artillery", Positions: []string{"Mid", "Support"}, DamageType: "Magic", Range: "Ranged"},
	"Ziggs":   {Name: "Ziggs", Class: "Mage", Subclass: "Artillery", Positions: []string{"Mid", "Bot"}, DamageType: "Magic", Range: "Ranged"},

	// Marksmen
	"Aphelios": {Name: "Aphelios", Class: "Marksman", Subclass: "", Positions: []string{"Bot"}, DamageType: "Physical", Range: "Ranged"},
	"Ashe":     {Name: "Ashe", Class: "Marksman", Subclass: "", Positions: []string{"Bot"}, DamageType: "Physical", Range: "Ranged"},
	"Caitlyn":  {Name: "Caitlyn", Class: "Marksman", Subclass: "", Positions: []string{"Bot"}, DamageType: "Physical", Range: "Ranged"},
	"Corki":    {Name: "Corki", Class: "Marksman", Subclass: "", Positions: []string{"Mid"}, DamageType: "Mixed", Range: "Ranged"},
	"Draven":   {Name: "Draven", Class: "Marksman", Subclass: "", Positions: []string{"Bot"}, DamageType: "Physical", Range: "Ranged"},
	"Ezreal":   {Name: "Ezreal", Class: "Marksman", Subclass: "", Positions: []string{"Bot"}, DamageType: "Mixed", Range: "Ranged"},
	"Jhin":     {Name: "Jhin", Class: "Marksman", Subclass: "", Positions: []string{"Bot"}, DamageType: "Physical", Range: "Ranged"},
	"Jinx":     {Name: "Jinx", Class: "Marksman", Subclass: "", Positions: []string{"Bot"}, DamageType: "Physical", Range: "Ranged"},
	"Kai'Sa":   {Name: "Kai'Sa", Class: "Marksman", Subclass: "", Positions: []string{"Bot"}, DamageType: "Mixed", Range: "Ranged"},
	"Kalista":  {Name: "Kalista", Class: "Marksman", Subclass: "", Positions: []string{"Bot"}, DamageType: "Physical", Range: "Ranged"},
	"Kindred":  {Name: "Kindred", Class: "Marksman", Subclass: "", Positions: []string{"Jungle"}, DamageType: "Physical", Range: "Ranged"},
	"Lucian":   {Name: "Lucian", Class: "Marksman", Subclass: "", Positions: []string{"Bot", "Mid"}, DamageType: "Physical", Range: "Ranged"},
	"Miss Fortune": {Name: "Miss Fortune", Class: "Marksman", Subclass: "", Positions: []string{"Bot"}, DamageType: "Physical", Range: "Ranged"},
	"Nilah":    {Name: "Nilah", Class: "Marksman", Subclass: "", Positions: []string{"Bot"}, DamageType: "Physical", Range: "Melee"},
	"Samira":   {Name: "Samira", Class: "Marksman", Subclass: "", Positions: []string{"Bot"}, DamageType: "Physical", Range: "Ranged"},
	"Senna":    {Name: "Senna", Class: "Marksman", Subclass: "", Positions: []string{"Support", "Bot"}, DamageType: "Physical", Range: "Ranged"},
	"Sivir":    {Name: "Sivir", Class: "Marksman", Subclass: "", Positions: []string{"Bot"}, DamageType: "Physical", Range: "Ranged"},
	"Smolder":  {Name: "Smolder", Class: "Marksman", Subclass: "", Positions: []string{"Bot", "Mid"}, DamageType: "Magic", Range: "Ranged"},
	"Tristana": {Name: "Tristana", Class: "Marksman", Subclass: "", Positions: []string{"Bot", "Mid"}, DamageType: "Physical", Range: "Ranged"},
	"Twitch":   {Name: "Twitch", Class: "Marksman", Subclass: "", Positions: []string{"Bot"}, DamageType: "Physical", Range: "Ranged"},
	"Varus":    {Name: "Varus", Class: "Marksman", Subclass: "", Positions: []string{"Bot"}, DamageType: "Mixed", Range: "Ranged"},
	"Vayne":    {Name: "Vayne", Class: "Marksman", Subclass: "", Positions: []string{"Bot", "Top"}, DamageType: "Physical", Range: "Ranged"},
	"Xayah":    {Name: "Xayah", Class: "Marksman", Subclass: "", Positions: []string{"Bot"}, DamageType: "Physical", Range: "Ranged"},
	"Zeri":     {Name: "Zeri", Class: "Marksman", Subclass: "", Positions: []string{"Bot"}, DamageType: "Physical", Range: "Ranged"},

	// Slayers - Assassins
	"Akali":    {Name: "Akali", Class: "Slayer", Subclass: "Assassin", Positions: []string{"Mid", "Top"}, DamageType: "Magic", Range: "Melee"},
	"Akshan":   {Name: "Akshan", Class: "Slayer", Subclass: "Assassin", Positions: []string{"Mid"}, DamageType: "Physical", Range: "Ranged"},
	"Ekko":     {Name: "Ekko", Class: "Slayer", Subclass: "Assassin", Positions: []string{"Mid", "Jungle"}, DamageType: "Magic", Range: "Melee"},
	"Evelynn":  {Name: "Evelynn", Class: "Slayer", Subclass: "Assassin", Positions: []string{"Jungle"}, DamageType: "Magic", Range: "Melee"},
	"Fizz":     {Name: "Fizz", Class: "Slayer", Subclass: "Assassin", Positions: []string{"Mid"}, DamageType: "Magic", Range: "Melee"},
	"Kassadin": {Name: "Kassadin", Class: "Slayer", Subclass: "Assassin", Positions: []string{"Mid"}, DamageType: "Magic", Range: "Melee"},
	"Katarina": {Name: "Katarina", Class: "Slayer", Subclass: "Assassin", Positions: []string{"Mid"}, DamageType: "Magic", Range: "Melee"},
	"Kha'Zix":  {Name: "Kha'Zix", Class: "Slayer", Subclass: "Assassin", Positions: []string{"Jungle"}, DamageType: "Physical", Range: "Melee"},
	"Naafiri":  {Name: "Naafiri", Class: "Slayer", Subclass: "Assassin", Positions: []string{"Mid", "Jungle"}, DamageType: "Physical", Range: "Melee"},
	"Nidalee":  {Name: "Nidalee", Class: "Slayer", Subclass: "Assassin", Positions: []string{"Jungle"}, DamageType: "Magic", Range: "Ranged"},
	"Nocturne": {Name: "Nocturne", Class: "Slayer", Subclass: "Assassin", Positions: []string{"Jungle"}, DamageType: "Physical", Range: "Melee"},
	"Qiyana":   {Name: "Qiyana", Class: "Slayer", Subclass: "Assassin", Positions: []string{"Mid", "Jungle"}, DamageType: "Physical", Range: "Melee"},
	"Rengar":   {Name: "Rengar", Class: "Slayer", Subclass: "Assassin", Positions: []string{"Jungle", "Top"}, DamageType: "Physical", Range: "Melee"},
	"Shaco":    {Name: "Shaco", Class: "Slayer", Subclass: "Assassin", Positions: []string{"Jungle"}, DamageType: "Physical", Range: "Melee"},
	"Talon":    {Name: "Talon", Class: "Slayer", Subclass: "Assassin", Positions: []string{"Mid", "Jungle"}, DamageType: "Physical", Range: "Melee"},
	"Zed":      {Name: "Zed", Class: "Slayer", Subclass: "Assassin", Positions: []string{"Mid"}, DamageType: "Physical", Range: "Melee"},

	// Slayers - Skirmishers
	"Bel'Veth":   {Name: "Bel'Veth", Class: "Slayer", Subclass: "Skirmisher", Positions: []string{"Jungle"}, DamageType: "Physical", Range: "Melee"},
	"Fiora":      {Name: "Fiora", Class: "Slayer", Subclass: "Skirmisher", Positions: []string{"Top"}, DamageType: "Physical", Range: "Melee"},
	"Gwen":       {Name: "Gwen", Class: "Slayer", Subclass: "Skirmisher", Positions: []string{"Top"}, DamageType: "Magic", Range: "Melee"},
	"Jax":        {Name: "Jax", Class: "Slayer", Subclass: "Skirmisher", Positions: []string{"Top", "Jungle"}, DamageType: "Mixed", Range: "Melee"},
	"Kayn":       {Name: "Kayn", Class: "Slayer", Subclass: "Skirmisher", Positions: []string{"Jungle"}, DamageType: "Physical", Range: "Melee"},
	"Lillia":     {Name: "Lillia", Class: "Slayer", Subclass: "Skirmisher", Positions: []string{"Jungle"}, DamageType: "Magic", Range: "Melee"},
	"Master Yi":  {Name: "Master Yi", Class: "Slayer", Subclass: "Skirmisher", Positions: []string{"Jungle"}, DamageType: "Physical", Range: "Melee"},
	"Riven":      {Name: "Riven", Class: "Slayer", Subclass: "Skirmisher", Positions: []string{"Top"}, DamageType: "Physical", Range: "Melee"},
	"Sylas":      {Name: "Sylas", Class: "Slayer", Subclass: "Skirmisher", Positions: []string{"Mid", "Jungle"}, DamageType: "Magic", Range: "Melee"},
	"Tryndamere": {Name: "Tryndamere", Class: "Slayer", Subclass: "Skirmisher", Positions: []string{"Top"}, DamageType: "Physical", Range: "Melee"},
	"Viego":      {Name: "Viego", Class: "Slayer", Subclass: "Skirmisher", Positions: []string{"Jungle"}, DamageType: "Physical", Range: "Melee"},
	"Yasuo":      {Name: "Yasuo", Class: "Slayer", Subclass: "Skirmisher", Positions: []string{"Mid", "Top", "Bot"}, DamageType: "Physical", Range: "Melee"},
	"Yone":       {Name: "Yone", Class: "Slayer", Subclass: "Skirmisher", Positions: []string{"Mid", "Top"}, DamageType: "Mixed", Range: "Melee"},

	// Tanks - Vanguards
	"Alistar":  {Name: "Alistar", Class: "Tank", Subclass: "Vanguard", Positions: []string{"Support"}, DamageType: "Magic", Range: "Melee"},
	"Amumu":    {Name: "Amumu", Class: "Tank", Subclass: "Vanguard", Positions: []string{"Jungle", "Support"}, DamageType: "Magic", Range: "Melee"},
	"Gragas":   {Name: "Gragas", Class: "Tank", Subclass: "Vanguard", Positions: []string{"Jungle", "Top"}, DamageType: "Magic", Range: "Melee"},
	"Leona":    {Name: "Leona", Class: "Tank", Subclass: "Vanguard", Positions: []string{"Support"}, DamageType: "Magic", Range: "Melee"},
	"Malphite": {Name: "Malphite", Class: "Tank", Subclass: "Vanguard", Positions: []string{"Top", "Support"}, DamageType: "Magic", Range: "Melee"},
	"Maokai":   {Name: "Maokai", Class: "Tank", Subclass: "Vanguard", Positions: []string{"Support", "Jungle"}, DamageType: "Magic", Range: "Melee"},
	"Nunu":     {Name: "Nunu & Willump", Class: "Tank", Subclass: "Vanguard", Positions: []string{"Jungle"}, DamageType: "Magic", Range: "Melee"},
	"Ornn":     {Name: "Ornn", Class: "Tank", Subclass: "Vanguard", Positions: []string{"Top"}, DamageType: "Magic", Range: "Melee"},
	"Rammus":   {Name: "Rammus", Class: "Tank", Subclass: "Vanguard", Positions: []string{"Jungle"}, DamageType: "Magic", Range: "Melee"},
	"Rell":     {Name: "Rell", Class: "Tank", Subclass: "Vanguard", Positions: []string{"Support"}, DamageType: "Magic", Range: "Melee"},
	"Sejuani":  {Name: "Sejuani", Class: "Tank", Subclass: "Vanguard", Positions: []string{"Jungle"}, DamageType: "Magic", Range: "Melee"},
	"Skarner":  {Name: "Skarner", Class: "Tank", Subclass: "Vanguard", Positions: []string{"Jungle"}, DamageType: "Mixed", Range: "Melee"},
	"Zac":      {Name: "Zac", Class: "Tank", Subclass: "Vanguard", Positions: []string{"Jungle"}, DamageType: "Magic", Range: "Melee"},

	// Tanks - Wardens
	"Braum":      {Name: "Braum", Class: "Tank", Subclass: "Warden", Positions: []string{"Support"}, DamageType: "Magic", Range: "Melee"},
	"Cho'Gath":   {Name: "Cho'Gath", Class: "Tank", Subclass: "Warden", Positions: []string{"Top"}, DamageType: "Magic", Range: "Melee"},
	"Galio":      {Name: "Galio", Class: "Tank", Subclass: "Warden", Positions: []string{"Mid", "Support"}, DamageType: "Magic", Range: "Melee"},
	"K'Sante":    {Name: "K'Sante", Class: "Tank", Subclass: "Warden", Positions: []string{"Top"}, DamageType: "Physical", Range: "Melee"},
	"Poppy":      {Name: "Poppy", Class: "Tank", Subclass: "Warden", Positions: []string{"Jungle", "Top"}, DamageType: "Physical", Range: "Melee"},
	"Shen":       {Name: "Shen", Class: "Tank", Subclass: "Warden", Positions: []string{"Top", "Support"}, DamageType: "Magic", Range: "Melee"},
	"Tahm Kench": {Name: "Tahm Kench", Class: "Tank", Subclass: "Warden", Positions: []string{"Support", "Top"}, DamageType: "Magic", Range: "Melee"},

	// Specialists
	"Azir":        {Name: "Azir", Class: "Specialist", Subclass: "", Positions: []string{"Mid"}, DamageType: "Magic", Range: "Ranged"},
	"Fiddlesticks": {Name: "Fiddlesticks", Class: "Specialist", Subclass: "", Positions: []string{"Jungle"}, DamageType: "Magic", Range: "Ranged"},
	"Gangplank":   {Name: "Gangplank", Class: "Specialist", Subclass: "", Positions: []string{"Top"}, DamageType: "Physical", Range: "Melee"},
	"Gnar":        {Name: "Gnar", Class: "Specialist", Subclass: "", Positions: []string{"Top"}, DamageType: "Physical", Range: "Ranged"},
	"Graves":      {Name: "Graves", Class: "Specialist", Subclass: "", Positions: []string{"Jungle"}, DamageType: "Physical", Range: "Ranged"},
	"Heimerdinger": {Name: "Heimerdinger", Class: "Specialist", Subclass: "", Positions: []string{"Top", "Mid", "Support"}, DamageType: "Magic", Range: "Ranged"},
	"Kayle":       {Name: "Kayle", Class: "Specialist", Subclass: "", Positions: []string{"Top"}, DamageType: "Magic", Range: "Melee"},
	"Kennen":      {Name: "Kennen", Class: "Specialist", Subclass: "", Positions: []string{"Top"}, DamageType: "Magic", Range: "Ranged"},
	"Quinn":       {Name: "Quinn", Class: "Specialist", Subclass: "", Positions: []string{"Top"}, DamageType: "Physical", Range: "Ranged"},
	"Singed":      {Name: "Singed", Class: "Specialist", Subclass: "", Positions: []string{"Top"}, DamageType: "Magic", Range: "Melee"},
	"Teemo":       {Name: "Teemo", Class: "Specialist", Subclass: "", Positions: []string{"Top"}, DamageType: "Magic", Range: "Ranged"},
	"Zilean":      {Name: "Zilean", Class: "Specialist", Subclass: "", Positions: []string{"Support", "Mid"}, DamageType: "Magic", Range: "Ranged"},
}


// AgentData contains metadata about a VALORANT agent
// Data sourced from official VALORANT game
type AgentData struct {
	Name        string   `json:"name"`
	Role        string   `json:"role"`        // Duelist, Initiator, Controller, Sentinel
	Abilities   []string `json:"abilities"`   // List of ability names
	PlayStyle   string   `json:"playStyle"`   // Aggressive, Supportive, Defensive, Versatile
	Difficulty  string   `json:"difficulty"`  // Easy, Medium, Hard
}

// AgentDatabase is the comprehensive database of all VALORANT agents with their roles
// Based on official VALORANT agent classifications
var AgentDatabase = map[string]AgentData{
	// Duelists - Entry fraggers, aggressive playstyle
	"Jett":    {Name: "Jett", Role: "Duelist", Abilities: []string{"Cloudburst", "Updraft", "Tailwind", "Blade Storm"}, PlayStyle: "Aggressive", Difficulty: "Medium"},
	"Phoenix": {Name: "Phoenix", Role: "Duelist", Abilities: []string{"Blaze", "Curveball", "Hot Hands", "Run It Back"}, PlayStyle: "Aggressive", Difficulty: "Easy"},
	"Reyna":   {Name: "Reyna", Role: "Duelist", Abilities: []string{"Leer", "Devour", "Dismiss", "Empress"}, PlayStyle: "Aggressive", Difficulty: "Easy"},
	"Raze":    {Name: "Raze", Role: "Duelist", Abilities: []string{"Boom Bot", "Blast Pack", "Paint Shells", "Showstopper"}, PlayStyle: "Aggressive", Difficulty: "Medium"},
	"Yoru":    {Name: "Yoru", Role: "Duelist", Abilities: []string{"Fakeout", "Blindside", "Gatecrash", "Dimensional Drift"}, PlayStyle: "Aggressive", Difficulty: "Hard"},
	"Neon":    {Name: "Neon", Role: "Duelist", Abilities: []string{"Fast Lane", "Relay Bolt", "High Gear", "Overdrive"}, PlayStyle: "Aggressive", Difficulty: "Hard"},
	"Iso":     {Name: "Iso", Role: "Duelist", Abilities: []string{"Undercut", "Double Tap", "Contingency", "Kill Contract"}, PlayStyle: "Aggressive", Difficulty: "Medium"},
	"Waylay":  {Name: "Waylay", Role: "Duelist", Abilities: []string{"Ricochet", "Flashpoint", "Ambush", "Crossfire"}, PlayStyle: "Aggressive", Difficulty: "Medium"},

	// Initiators - Info gathering, entry support
	"Sova":   {Name: "Sova", Role: "Initiator", Abilities: []string{"Owl Drone", "Shock Bolt", "Recon Bolt", "Hunter's Fury"}, PlayStyle: "Supportive", Difficulty: "Medium"},
	"Breach": {Name: "Breach", Role: "Initiator", Abilities: []string{"Aftershock", "Flashpoint", "Fault Line", "Rolling Thunder"}, PlayStyle: "Aggressive", Difficulty: "Medium"},
	"Skye":   {Name: "Skye", Role: "Initiator", Abilities: []string{"Regrowth", "Trailblazer", "Guiding Light", "Seekers"}, PlayStyle: "Supportive", Difficulty: "Medium"},
	"KAY/O":  {Name: "KAY/O", Role: "Initiator", Abilities: []string{"FRAG/ment", "FLASH/drive", "ZERO/point", "NULL/cmd"}, PlayStyle: "Aggressive", Difficulty: "Easy"},
	"Fade":   {Name: "Fade", Role: "Initiator", Abilities: []string{"Prowler", "Seize", "Haunt", "Nightfall"}, PlayStyle: "Supportive", Difficulty: "Medium"},
	"Gekko":  {Name: "Gekko", Role: "Initiator", Abilities: []string{"Mosh Pit", "Wingman", "Dizzy", "Thrash"}, PlayStyle: "Supportive", Difficulty: "Easy"},
	"Tejo":   {Name: "Tejo", Role: "Initiator", Abilities: []string{"Stealth Drone", "Guided Salvo", "Special Delivery", "Armageddon"}, PlayStyle: "Aggressive", Difficulty: "Medium"},

	// Controllers - Smoke, area denial
	"Brimstone": {Name: "Brimstone", Role: "Controller", Abilities: []string{"Stim Beacon", "Incendiary", "Sky Smoke", "Orbital Strike"}, PlayStyle: "Supportive", Difficulty: "Easy"},
	"Omen":      {Name: "Omen", Role: "Controller", Abilities: []string{"Shrouded Step", "Paranoia", "Dark Cover", "From the Shadows"}, PlayStyle: "Versatile", Difficulty: "Medium"},
	"Viper":     {Name: "Viper", Role: "Controller", Abilities: []string{"Snake Bite", "Poison Cloud", "Toxic Screen", "Viper's Pit"}, PlayStyle: "Defensive", Difficulty: "Hard"},
	"Astra":     {Name: "Astra", Role: "Controller", Abilities: []string{"Gravity Well", "Nova Pulse", "Nebula", "Cosmic Divide"}, PlayStyle: "Supportive", Difficulty: "Hard"},
	"Harbor":    {Name: "Harbor", Role: "Controller", Abilities: []string{"Cascade", "Cove", "High Tide", "Reckoning"}, PlayStyle: "Supportive", Difficulty: "Medium"},
	"Clove":     {Name: "Clove", Role: "Controller", Abilities: []string{"Pick-Me-Up", "Meddle", "Ruse", "Not Dead Yet"}, PlayStyle: "Aggressive", Difficulty: "Medium"},

	// Sentinels - Defense, site anchor
	"Sage":     {Name: "Sage", Role: "Sentinel", Abilities: []string{"Barrier Orb", "Slow Orb", "Healing Orb", "Resurrection"}, PlayStyle: "Defensive", Difficulty: "Easy"},
	"Cypher":   {Name: "Cypher", Role: "Sentinel", Abilities: []string{"Trapwire", "Cyber Cage", "Spycam", "Neural Theft"}, PlayStyle: "Defensive", Difficulty: "Medium"},
	"Killjoy":  {Name: "Killjoy", Role: "Sentinel", Abilities: []string{"Nanoswarm", "Alarmbot", "Turret", "Lockdown"}, PlayStyle: "Defensive", Difficulty: "Easy"},
	"Chamber":  {Name: "Chamber", Role: "Sentinel", Abilities: []string{"Trademark", "Headhunter", "Rendezvous", "Tour De Force"}, PlayStyle: "Aggressive", Difficulty: "Medium"},
	"Deadlock": {Name: "Deadlock", Role: "Sentinel", Abilities: []string{"GravNet", "Sonic Sensor", "Barrier Mesh", "Annihilation"}, PlayStyle: "Defensive", Difficulty: "Medium"},
	"Vyse":     {Name: "Vyse", Role: "Sentinel", Abilities: []string{"Shear", "Arc Rose", "Razorvine", "Steel Garden"}, PlayStyle: "Defensive", Difficulty: "Medium"},
}

// GetChampionClass returns the class of a champion
func GetChampionClass(championName string) string {
	if data, exists := ChampionDatabase[championName]; exists {
		return data.Class
	}
	return "Unknown"
}

// GetChampionSubclass returns the subclass of a champion
func GetChampionSubclass(championName string) string {
	if data, exists := ChampionDatabase[championName]; exists {
		return data.Subclass
	}
	return "Unknown"
}

// GetChampionData returns full champion data
func GetChampionData(championName string) (ChampionData, bool) {
	data, exists := ChampionDatabase[championName]
	return data, exists
}

// GetAgentRole returns the role of a VALORANT agent
func GetAgentRole(agentName string) string {
	if data, exists := AgentDatabase[agentName]; exists {
		return data.Role
	}
	return "Unknown"
}

// GetChampionRole returns the primary role/position of a LoL champion
func GetChampionRole(championName string) string {
	if data, exists := ChampionDatabase[championName]; exists {
		if len(data.Positions) > 0 {
			return data.Positions[0] // Return primary position
		}
	}
	return ""
}

// GetAgentData returns full agent data
func GetAgentData(agentName string) (AgentData, bool) {
	data, exists := AgentDatabase[agentName]
	return data, exists
}

// ClassifyLoLCompositionArchetype analyzes a team composition and returns the archetype
// Uses real champion data to determine composition style
func ClassifyLoLCompositionArchetype(champions []string) string {
	if len(champions) == 0 {
		return "unknown"
	}

	// Count classes and subclasses
	classCounts := make(map[string]int)
	subclassCounts := make(map[string]int)
	damageTypes := make(map[string]int)

	for _, champ := range champions {
		if data, exists := ChampionDatabase[champ]; exists {
			classCounts[data.Class]++
			if data.Subclass != "" {
				subclassCounts[data.Subclass]++
			}
			damageTypes[data.DamageType]++
		}
	}

	// Determine archetype based on composition
	// Teamfight: Multiple Vanguards, Battlemages, or AoE-heavy champions
	teamfightScore := subclassCounts["Vanguard"] + subclassCounts["Battlemage"] + subclassCounts["Burst"]
	
	// Pick: Multiple Assassins, Catchers, or single-target burst
	pickScore := subclassCounts["Assassin"] + subclassCounts["Catcher"] + subclassCounts["Burst"]
	
	// Split-push: Skirmishers, Juggernauts with dueling power
	splitScore := subclassCounts["Skirmisher"] + subclassCounts["Juggernaut"]
	
	// Scaling: Multiple Marksmen, late-game carries
	scaleScore := classCounts["Marksman"]
	if subclassCounts["Skirmisher"] > 0 {
		scaleScore++ // Skirmishers often scale well
	}

	// Poke: Artillery mages, long-range champions
	pokeScore := subclassCounts["Artillery"]

	// Dive: Divers, Assassins
	diveScore := subclassCounts["Diver"] + subclassCounts["Assassin"]

	// Protect the carry: Enchanters + Marksmen
	protectScore := 0
	if classCounts["Marksman"] >= 1 && subclassCounts["Enchanter"] >= 1 {
		protectScore = classCounts["Marksman"] + subclassCounts["Enchanter"] + subclassCounts["Warden"]
	}

	// Find the dominant archetype
	maxScore := teamfightScore
	archetype := "teamfight"

	if pickScore > maxScore {
		maxScore = pickScore
		archetype = "pick"
	}
	if splitScore > maxScore {
		maxScore = splitScore
		archetype = "split-push"
	}
	if scaleScore > maxScore {
		maxScore = scaleScore
		archetype = "scaling"
	}
	if pokeScore > maxScore {
		maxScore = pokeScore
		archetype = "poke"
	}
	if diveScore > maxScore {
		maxScore = diveScore
		archetype = "dive"
	}
	if protectScore > maxScore {
		maxScore = protectScore
		archetype = "protect-the-carry"
	}

	if maxScore == 0 {
		return "balanced"
	}

	return archetype
}

// ClassifyVALCompositionStyle analyzes a VALORANT team composition
func ClassifyVALCompositionStyle(agents []string) string {
	if len(agents) == 0 {
		return "unknown"
	}

	roleCounts := make(map[string]int)
	for _, agent := range agents {
		if data, exists := AgentDatabase[agent]; exists {
			roleCounts[data.Role]++
		}
	}

	// Standard comp: 1 Duelist, 1-2 Initiators, 1 Controller, 1-2 Sentinels
	duelists := roleCounts["Duelist"]
	initiators := roleCounts["Initiator"]
	controllers := roleCounts["Controller"]
	sentinels := roleCounts["Sentinel"]

	// Aggressive: 2+ Duelists
	if duelists >= 2 {
		return "aggressive"
	}

	// Defensive: 2+ Sentinels
	if sentinels >= 2 {
		return "defensive"
	}

	// Execute-heavy: 2+ Initiators
	if initiators >= 2 {
		return "execute-heavy"
	}

	// Control-heavy: 2+ Controllers
	if controllers >= 2 {
		return "control-heavy"
	}

	// Standard balanced comp
	if duelists >= 1 && controllers >= 1 && (initiators >= 1 || sentinels >= 1) {
		return "standard"
	}

	return "unconventional"
}


// =============================================================================
// HACKATHON-WINNING: Class Category Functions for Matchup Analysis
// =============================================================================

// GetClassCategory returns a simplified class category for matchup analysis
// Categories: "assassin", "mage", "tank", "fighter", "marksman", "support"
// This enables insights like "2.1 KDA on control mages vs 8.5 on assassins"
func GetClassCategory(championName string) string {
	data, exists := ChampionDatabase[championName]
	if !exists {
		return "unknown"
	}

	// Map Riot's official classes to simplified categories
	switch data.Class {
	case "Slayer":
		// Slayers include Assassins and Skirmishers
		if data.Subclass == "Assassin" {
			return "assassin"
		}
		return "fighter" // Skirmishers are fighter-like
	case "Mage":
		return "mage"
	case "Tank":
		return "tank"
	case "Fighter":
		return "fighter"
	case "Marksman":
		return "marksman"
	case "Controller":
		return "support"
	case "Specialist":
		// Specialists are varied - classify by subclass or damage type
		if data.DamageType == "Magic" {
			return "mage"
		}
		return "fighter"
	default:
		return "unknown"
	}
}

// GetMatchupModifier returns the expected performance modifier for a class matchup
// Positive = favorable, Negative = unfavorable
// Used for identifying exploitable matchups
func GetMatchupModifier(playerClass, opponentClass string) float64 {
	modifiers := map[string]map[string]float64{
		"assassin": {
			"mage":     0.15,  // Assassins generally beat mages
			"marksman": 0.20,  // Assassins excel at killing ADCs
			"tank":     -0.15, // Assassins struggle vs tanks
			"fighter":  -0.10, // Fighters can duel assassins
			"support":  0.05,  // Neutral to slight advantage
		},
		"mage": {
			"assassin": -0.15, // Mages vulnerable to assassins
			"tank":     0.10,  // Mages can kite tanks
			"fighter":  0.05,  // Slight advantage with range
			"marksman": 0.0,   // Neutral
			"support":  0.0,   // Neutral
		},
		"tank": {
			"assassin": 0.15,  // Tanks survive assassin burst
			"mage":     -0.10, // Mages can kite
			"fighter":  -0.05, // Fighters can duel
			"marksman": -0.10, // ADCs shred tanks late
			"support":  0.10,  // Tanks beat supports in lane
		},
		"fighter": {
			"assassin": 0.10,  // Fighters can duel
			"mage":     -0.05, // Mages have range
			"tank":     0.05,  // Fighters beat tanks 1v1
			"marksman": 0.15,  // Fighters dive ADCs
			"support":  0.10,  // Fighters beat supports
		},
		"marksman": {
			"assassin": -0.20, // ADCs die to assassins
			"mage":     0.0,   // Neutral
			"tank":     0.10,  // ADCs shred tanks late
			"fighter":  -0.15, // Fighters dive ADCs
			"support":  0.05,  // ADCs beat supports
		},
	}

	if playerMods, ok := modifiers[playerClass]; ok {
		if mod, ok := playerMods[opponentClass]; ok {
			return mod
		}
	}
	return 0.0
}

// IsAssassin returns true if the champion is an assassin
func IsAssassin(championName string) bool {
	data, exists := ChampionDatabase[championName]
	if !exists {
		return false
	}
	return data.Subclass == "Assassin"
}

// IsMage returns true if the champion is a mage (any type)
func IsMage(championName string) bool {
	data, exists := ChampionDatabase[championName]
	if !exists {
		return false
	}
	return data.Class == "Mage"
}

// IsTank returns true if the champion is a tank
func IsTank(championName string) bool {
	data, exists := ChampionDatabase[championName]
	if !exists {
		return false
	}
	return data.Class == "Tank"
}

// IsFighter returns true if the champion is a fighter/bruiser
func IsFighter(championName string) bool {
	data, exists := ChampionDatabase[championName]
	if !exists {
		return false
	}
	return data.Class == "Fighter" || data.Subclass == "Skirmisher"
}

// IsMarksman returns true if the champion is a marksman/ADC
func IsMarksman(championName string) bool {
	data, exists := ChampionDatabase[championName]
	if !exists {
		return false
	}
	return data.Class == "Marksman"
}

// IsSupport returns true if the champion is typically played as support
func IsSupport(championName string) bool {
	data, exists := ChampionDatabase[championName]
	if !exists {
		return false
	}
	return data.Class == "Controller"
}
