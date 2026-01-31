import type { GameTitle, Team } from '@/types';

// Comprehensive team lists from GRID API
// These are pre-loaded to avoid API calls for team search
// Last updated: January 2026

export const ALL_TEAMS: Record<GameTitle, Team[]> = {
  lol: [
    // LCK (Korea)
    { id: '47494', name: 'T1', logoUrl: 'https://cdn.grid.gg/assets/team-logos/1e7311945adc58ac807ffcf10b18d002' },
    { id: '47558', name: 'Gen.G Esports', logoUrl: 'https://cdn.grid.gg/assets/team-logos/d2eded1af01ce76afb9540de0ef8b1d8' },
    { id: '406', name: 'Hanwha Life Esports', logoUrl: 'https://cdn.grid.gg/assets/team-logos/f6bbce9ba43dfbf1b50b6cde51fda71b' },
    { id: '47961', name: 'DRX', logoUrl: 'https://cdn.grid.gg/assets/team-logos/6470bf630495e659e6120d516a2f790c' },
    { id: '48179', name: 'Dplus KIA', logoUrl: 'https://cdn.grid.gg/assets/team-logos/1c4e991b3a2ec38bc188409b6dcf6427' },
    { id: '407', name: 'KT Rolster', logoUrl: 'https://cdn.grid.gg/assets/team-logos/a47aaabd94d8ee66fc22a6893a48f4ae' },
    { id: '52747', name: 'Nongshim RedForce', logoUrl: 'https://cdn.grid.gg/assets/team-logos/15cf94cff3b13fd908e2b79576b8e6f0' },
    { id: '4035', name: 'BNK FearX', logoUrl: 'https://cdn.grid.gg/assets/team-logos/0ee8dc4cac1c6b09c4b25b7cbefc2493' },
    { id: '52817', name: 'BRION', logoUrl: 'https://cdn.grid.gg/assets/team-logos/8fcd5ca8c455d1173afc0815b7321b7a' },
    { id: '3483', name: 'DN SOOPers', logoUrl: 'https://cdn.grid.gg/assets/team-logos/cfcf51afe7a2d0b7559a391e9b39f8e6' },
    
    // LEC (Europe)
    { id: '47376', name: 'Fnatic', logoUrl: 'https://cdn.grid.gg/assets/team-logos/d5bd0cb8ca32672cd8608d2ad2cb039a' },
    { id: '47380', name: 'G2 Esports', logoUrl: 'https://cdn.grid.gg/assets/team-logos/94a774753c28adb6602d3c36e428a849' },
    { id: '53168', name: 'GIANTX', logoUrl: 'https://cdn.grid.gg/assets/team-logos/af80b36401e45be0a5eb13549ff4128f' },
    { id: '53165', name: 'Karmine Corp', logoUrl: 'https://cdn.grid.gg/assets/team-logos/0c73991760c7d80df981a06e99c7cd51' },
    { id: '47619', name: 'Movistar KOI', logoUrl: 'https://cdn.grid.gg/assets/team-logos/739e7667002e441bd2596ba7fecff107' },
    { id: '106', name: 'Rogue', logoUrl: 'https://cdn.grid.gg/assets/team-logos/9089f72e65b99e4b08e88815fffa5576' },
    { id: '353', name: 'SK Gaming', logoUrl: 'https://cdn.grid.gg/assets/team-logos/c69f9b19774472b7b06d183c9de2a448' },
    { id: '47435', name: 'Team Heretics', logoUrl: 'https://cdn.grid.gg/assets/team-logos/64b4c8501ef44ba51648a38d131d9f2e' },
    { id: '47370', name: 'Team Vitality', logoUrl: 'https://cdn.grid.gg/assets/team-logos/195171114bca3c1f69136199f73fce14' },
    { id: '52661', name: 'Shifters', logoUrl: 'https://cdn.grid.gg/assets/team-logos/7df957bb35378ac0ba4dcd1cdb50ee83' },
    { id: '346', name: 'Excel Esports', logoUrl: 'https://cdn.grid.gg/assets/team-logos/1daee6ec4b66dcfc37d2c5525492bc34' },
    { id: '348', name: 'Misfits Gaming', logoUrl: 'https://cdn.grid.gg/assets/team-logos/8ec6989cd15d1201922fd4fa65a83011' },
    
    // LCS (North America)
    { id: '338', name: 'Counter Logic Gaming', logoUrl: 'https://cdn.grid.gg/assets/team-logos/204579fb71736a1228c3bf63a01de743' },
    { id: '340', name: 'FlyQuest', logoUrl: 'https://cdn.grid.gg/assets/team-logos/c073b6e06dbed1d34d37fdcfdda85e4d' },
    { id: '343', name: 'Golden Guardians', logoUrl: 'https://cdn.grid.gg/assets/team-logos/bc271cd59aab45cae566ac1feb7a2c67' },
    { id: '345', name: 'TSM', logoUrl: 'https://cdn.grid.gg/assets/team-logos/ee960ebd08b9d1fcfc01b4356c7b31fc' },
    { id: '344', name: 'Dignitas', logoUrl: 'https://cdn.grid.gg/assets/team-logos/8ace50d131dda8bd071906f30e873e7d' },
    { id: '351', name: 'Immortals', logoUrl: 'https://cdn.grid.gg/assets/team-logos/f457570556aeddc7817a2546ef721d6e' },
    { id: '79', name: 'Cloud9', logoUrl: 'https://cdn.grid.gg/assets/team-logos/62118a302ef0900c7e7cc7b52eb2ca49' },
    { id: '83', name: 'Team Liquid', logoUrl: 'https://cdn.grid.gg/assets/team-logos/f70296f9d298ca02968da2948a7853ac' },
    { id: '337', name: '100 Thieves', logoUrl: 'https://cdn.grid.gg/assets/team-logos/d2a812773eb4be16fb663ab506843a2d' },
    { id: '281', name: 'Evil Geniuses', logoUrl: 'https://cdn.grid.gg/assets/team-logos/bc79e0a0bfed60984ca9a6a2005ff32e' },
    { id: '97', name: 'NRG', logoUrl: 'https://cdn.grid.gg/assets/team-logos/edc7c3f60522d6700ccc6afa083c3707' },
    { id: '99', name: 'FURIA', logoUrl: 'https://cdn.grid.gg/assets/team-logos/c5c3bc875f8e10b56ce095f0971c03f4' },
    
    // LPL (China)
    { id: '364', name: 'JD Gaming', logoUrl: 'https://cdn.grid.gg/assets/team-logos/12328ba2062d446064aeff9185688f6d' },
    { id: '366', name: 'LNG Esports', logoUrl: 'https://cdn.grid.gg/assets/team-logos/d7729d91b045a70632e3c64716f969d4' },
    { id: '375', name: 'Top Esports', logoUrl: 'https://cdn.grid.gg/assets/team-logos/0a3a4f1e3b62cc86bb86284a63e62521' },
    { id: '368', name: 'LGD Gaming', logoUrl: 'https://cdn.grid.gg/assets/team-logos/e6a643a21b7b0ac2253ff218993ee2da' },
    { id: '378', name: 'Rogue Warriors', logoUrl: 'https://cdn.grid.gg/assets/team-logos/eb10a3e26cbd2c4e2cef20aa23102e3c' },
    { id: '380', name: 'Victory Five', logoUrl: 'https://cdn.grid.gg/assets/team-logos/d9cff5e124dc669feeaf182643bc44e0' },
    { id: '373', name: 'Team WE', logoUrl: 'https://cdn.grid.gg/assets/team-logos/c4ed20dde97cd05bf43c9a1917a237ac' },
    { id: '369', name: 'Oh My God', logoUrl: 'https://cdn.grid.gg/assets/team-logos/a523933dbb1d3d560def090b6276fcf8' },
    { id: '359', name: 'EDward Gaming', logoUrl: 'https://cdn.grid.gg/assets/team-logos/f1f485acb151bd168248160d96e4e53b' },
    { id: '356', name: 'Bilibili Gaming', logoUrl: 'https://cdn.grid.gg/assets/team-logos/51307a108b3e4fdc5209079079a0ed37' },
    { id: '361', name: 'FunPlus Phoenix', logoUrl: 'https://cdn.grid.gg/assets/team-logos/90b55e249fbff2884884bef3b73ff054' },
    { id: '377', name: 'Vici Gaming', logoUrl: 'https://cdn.grid.gg/assets/team-logos/b0916afec11a90170bccc1052366ea96' },
    
    // Other regions
    { id: '418', name: 'GAM Esports', logoUrl: 'https://cdn.grid.gg/assets/team-logos/2e1619c0d54ae6fa6f9f48dce1407557' },
    { id: '400', name: 'Afreeca Freecs', logoUrl: 'https://cdn.grid.gg/assets/team-logos/97a76709f3b2d9616464de29829f0d65' },
    { id: '402', name: 'DAMWON Gaming', logoUrl: 'https://cdn.grid.gg/assets/team-logos/311d979499cc9807a0f59c01829467e0' },
    { id: '405', name: 'Griffin', logoUrl: 'https://cdn.grid.gg/assets/team-logos/c6dac34a2265c533e3fb421e77bcc2ac' },
    { id: '408', name: 'SANDBOX Gaming', logoUrl: 'https://cdn.grid.gg/assets/team-logos/c6a02983391dce28b952bfd4303e1912' },
    { id: '82', name: 'Natus Vincere', logoUrl: 'https://cdn.grid.gg/assets/team-logos/64a155fa2039125daee99c371ef478e8' },
    { id: '90', name: 'Team Vitality', logoUrl: 'https://cdn.grid.gg/assets/team-logos/d386b225a0563e211e798a6fb02c0e15' },
    { id: '94', name: 'Fnatic', logoUrl: 'https://cdn.grid.gg/assets/team-logos/fb84562e92e5ee0b0f66daba65c6ee65' },
    { id: '96', name: 'G2 Esports', logoUrl: 'https://cdn.grid.gg/assets/team-logos/f1735e5f2227b3ebc6e2bd1a38b9b092' },
    { id: '80', name: 'FaZe Clan', logoUrl: 'https://cdn.grid.gg/assets/team-logos/8219c228bea68e3798e170d7308ad76d' },
  ],
  valorant: [
    // VCT Americas
    { id: '1079', name: 'Sentinels', logoUrl: 'https://cdn.grid.gg/assets/team-logos/d05bdadd653b68669ad96d91d7f86cff' },
    { id: '79', name: 'Cloud9', logoUrl: 'https://cdn.grid.gg/assets/team-logos/62118a302ef0900c7e7cc7b52eb2ca49' },
    { id: '97', name: 'NRG', logoUrl: 'https://cdn.grid.gg/assets/team-logos/edc7c3f60522d6700ccc6afa083c3707' },
    { id: '3412', name: 'LOUD', logoUrl: 'https://cdn.grid.gg/assets/team-logos/e8d9a0765fab88054c132856085824be' },
    { id: '337', name: '100 Thieves', logoUrl: 'https://cdn.grid.gg/assets/team-logos/d2a812773eb4be16fb663ab506843a2d' },
    { id: '281', name: 'Evil Geniuses', logoUrl: 'https://cdn.grid.gg/assets/team-logos/bc79e0a0bfed60984ca9a6a2005ff32e' },
    { id: '99', name: 'FURIA', logoUrl: 'https://cdn.grid.gg/assets/team-logos/c5c3bc875f8e10b56ce095f0971c03f4' },
    { id: '81', name: 'MIBR', logoUrl: 'https://cdn.grid.gg/assets/team-logos/232c1e60a45d9f47bb5425a5ead72ff2' },
    { id: '1611', name: 'Leviatán', logoUrl: 'https://cdn.grid.gg/assets/team-logos/54f4e80969726869b0fbf48fb9ec5b18' },
    { id: '48457', name: 'KRÜ Esports', logoUrl: 'https://cdn.grid.gg/assets/team-logos/110044bfa66de1f7ba0672f13361a4db' },
    
    // VCT EMEA
    { id: '96', name: 'G2 Esports', logoUrl: 'https://cdn.grid.gg/assets/team-logos/f1735e5f2227b3ebc6e2bd1a38b9b092' },
    { id: '94', name: 'Fnatic', logoUrl: 'https://cdn.grid.gg/assets/team-logos/fb84562e92e5ee0b0f66daba65c6ee65' },
    { id: '173', name: 'Team Heretics', logoUrl: 'https://cdn.grid.gg/assets/team-logos/3fcbbe4a6ddb235481c105bc066831c6' },
    { id: '3554', name: 'Karmine Corp', logoUrl: 'https://cdn.grid.gg/assets/team-logos/03ff74c930bfb9ad4385287baac42e98' },
    { id: '83', name: 'Team Liquid', logoUrl: 'https://cdn.grid.gg/assets/team-logos/f70296f9d298ca02968da2948a7853ac' },
    { id: '90', name: 'Team Vitality', logoUrl: 'https://cdn.grid.gg/assets/team-logos/d386b225a0563e211e798a6fb02c0e15' },
    { id: '82', name: 'Natus Vincere', logoUrl: 'https://cdn.grid.gg/assets/team-logos/64a155fa2039125daee99c371ef478e8' },
    { id: '149', name: 'BIG', logoUrl: 'https://cdn.grid.gg/assets/team-logos/410a53c16616f3ce1c99f8007a0eea6c' },
    { id: '24', name: 'Ninjas in Pyjamas', logoUrl: 'https://cdn.grid.gg/assets/team-logos/35c97a5b06bac7e16406401265c97d42' },
    { id: '106', name: 'Rogue', logoUrl: 'https://cdn.grid.gg/assets/team-logos/9089f72e65b99e4b08e88815fffa5576' },
    
    // VCT Pacific
    { id: '917', name: 'Paper Rex', logoUrl: 'https://cdn.grid.gg/assets/team-logos/d17d7871f6b71772ffd62394e3966979' },
    { id: '359', name: 'EDward Gaming', logoUrl: 'https://cdn.grid.gg/assets/team-logos/f1f485acb151bd168248160d96e4e53b' },
    { id: '361', name: 'FunPlus Phoenix', logoUrl: 'https://cdn.grid.gg/assets/team-logos/90b55e249fbff2884884bef3b73ff054' },
    { id: '404', name: 'Gen.G', logoUrl: 'https://cdn.grid.gg/assets/team-logos/2734dd1c5942f3d8a115ff0d0360f0ad' },
    { id: '335', name: 'T1', logoUrl: 'https://cdn.grid.gg/assets/team-logos/f196184453d8f7874dc33621fd4d9807' },
    { id: '42', name: 'BOOM Esports', logoUrl: 'https://cdn.grid.gg/assets/team-logos/528a28f13f5756cfbfeac90ac78dd7a0' },
    { id: '407', name: 'DRX', logoUrl: 'https://cdn.grid.gg/assets/team-logos/a47aaabd94d8ee66fc22a6893a48f4ae' },
    { id: '169', name: 'BREN Esports', logoUrl: 'https://cdn.grid.gg/assets/team-logos/2d7f11ff64516c83ffaff2997d3d3256' },
    { id: '331', name: 'ORDER', logoUrl: 'https://cdn.grid.gg/assets/team-logos/9f281941b89ccb54c67ae312614f50a6' },
    
    // VCT China
    { id: '359', name: 'EDG', logoUrl: 'https://cdn.grid.gg/assets/team-logos/f1f485acb151bd168248160d96e4e53b' },
    { id: '361', name: 'FPX', logoUrl: 'https://cdn.grid.gg/assets/team-logos/90b55e249fbff2884884bef3b73ff054' },
    { id: '356', name: 'Bilibili Gaming', logoUrl: 'https://cdn.grid.gg/assets/team-logos/51307a108b3e4fdc5209079079a0ed37' },
    
    // Other notable teams
    { id: '80', name: 'FaZe Clan', logoUrl: 'https://cdn.grid.gg/assets/team-logos/8219c228bea68e3798e170d7308ad76d' },
    { id: '95', name: 'OpTic Gaming', logoUrl: 'https://cdn.grid.gg/assets/team-logos/006fcfa6617a77ed46b715ba0d640cde' },
    { id: '100', name: 'Luminosity Gaming', logoUrl: 'https://cdn.grid.gg/assets/team-logos/2d81367ef9755da16a8c5eaf667f32d8' },
    { id: '25', name: 'Alliance', logoUrl: 'https://cdn.grid.gg/assets/team-logos/5c50ebdce4b1890c640b663d9b4328a1' },
    { id: '29', name: 'Gambit Esports', logoUrl: 'https://cdn.grid.gg/assets/team-logos/265d4b8b7f01f1caf808b4e7a834faaa' },
    { id: '91', name: 'MOUZ', logoUrl: 'https://cdn.grid.gg/assets/team-logos/7b59d90baa74c1cc70598b6852001627' },
    { id: '69', name: 'Movistar Riders', logoUrl: 'https://cdn.grid.gg/assets/team-logos/c485a285bb900dcdac6a4798004c1d76' },
  ],
};

// Featured teams to display by default (subset of all teams)
export const FEATURED_TEAMS: Record<GameTitle, Team[]> = {
  lol: ALL_TEAMS.lol.slice(0, 12), // First 12 teams (top LCK + some LEC)
  valorant: ALL_TEAMS.valorant.slice(0, 12), // First 12 teams (top VCT teams)
};

// Helper function to search teams locally
export function searchTeamsLocally(query: string, game: GameTitle): Team[] {
  if (!query.trim()) return [];
  
  const lowerQuery = query.toLowerCase();
  return ALL_TEAMS[game].filter(team => 
    team.name.toLowerCase().includes(lowerQuery)
  );
}

// Helper function to get all teams for a game
export function getAllTeams(game: GameTitle): Team[] {
  return ALL_TEAMS[game];
}

// Helper function to find a team by ID
export function findTeamById(id: string, game: GameTitle): Team | undefined {
  return ALL_TEAMS[game].find(team => team.id === id);
}
