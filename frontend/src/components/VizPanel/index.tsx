import { t } from '../../i18n';
import { useSettingsStore, currentLang } from '../../store/settingsStore';
import type { VizModel } from '../../types';
import { ArrayViz } from './ArrayViz';
import { GraphViz } from './GraphViz';

export function VizPanel({ models }: { models: VizModel[] }) {
  const { settings } = useSettingsStore();
  const tr = t(currentLang(settings));

  if (!models.length) {
    return (
      <div className="flex h-full items-center justify-center p-4 text-xs text-muted">
        {tr.vizEmpty}
      </div>
    );
  }

  return (
    <div className="h-full overflow-y-auto">
      {models.map((m) => {
        if (m.type === 'array') {
          return <ArrayViz key={m.name} model={m} />;
        }
        if (m.type === 'linked_list' || m.type === 'tree' || m.type === 'map') {
          return <GraphViz key={m.name} model={m} />;
        }
        return null;
      })}
    </div>
  );
}
