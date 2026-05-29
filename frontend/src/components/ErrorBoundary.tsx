import { Component, type ErrorInfo, type ReactNode } from 'react';

interface Props {
  children: ReactNode;
}

interface State {
  error: string | null;
}

export class ErrorBoundary extends Component<Props, State> {
  state: State = { error: null };

  static getDerivedStateFromError(error: Error): State {
    return { error: error.message };
  }

  componentDidCatch(error: Error, info: ErrorInfo) {
    console.error('UI error:', error, info);
  }

  render() {
    if (this.state.error) {
      return (
        <div className="flex h-screen flex-col items-center justify-center bg-bg p-8 text-center">
          <h1 className="mb-2 text-lg font-semibold text-error">Ошибка интерфейса</h1>
          <p className="max-w-lg font-mono text-sm text-muted">{this.state.error}</p>
          <p className="mt-4 text-xs text-muted">Проверьте data/app.log рядом с visualisator.exe</p>
        </div>
      );
    }
    return this.props.children;
  }
}
