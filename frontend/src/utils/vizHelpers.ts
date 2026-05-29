import type { Variable } from '../types';

export function hasSliceLikeVariables(variables: Variable[]): boolean {
  for (const v of variables) {
    const typeLower = v.type.toLowerCase();
    const valLower = v.value.toLowerCase();
    if (
      typeLower.startsWith('[]') ||
      typeLower.startsWith('[') ||
      typeLower.startsWith('map[') ||
      typeLower.includes('slice') ||
      valLower.includes('len:')
    ) {
      return true;
    }
    if (v.children?.some((c) => c.name === 'Val' || c.name === 'Next' || c.name === 'Left')) {
      return true;
    }
  }
  return false;
}
