import * as migration_20251023_155957 from './20251023_155957';
import * as migration_20251023_183345 from './20251023_183345';

export const migrations = [
  {
    up: migration_20251023_155957.up,
    down: migration_20251023_155957.down,
    name: '20251023_155957',
  },
  {
    up: migration_20251023_183345.up,
    down: migration_20251023_183345.down,
    name: '20251023_183345'
  },
];
