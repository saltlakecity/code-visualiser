import type { VizModel } from '../../types';
import { ArrayViz } from './ArrayViz';

export function StringViz({ model }: { model: VizModel }) {
  return <ArrayViz model={{ ...model, name: `${model.name} (string)` }} />;
}
