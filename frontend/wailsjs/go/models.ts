export namespace models {
	
	export class OutlinePoint {
	    line: number;
	    text: string;
	
	    static createFrom(source: any = {}) {
	        return new OutlinePoint(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.line = source["line"];
	        this.text = source["text"];
	    }
	}
	export class SourceOutline {
	    name: string;
	    params: string;
	    returns: string;
	    loops?: OutlinePoint[];
	    branches?: OutlinePoint[];
	    returnsAt?: OutlinePoint[];
	    demoCall?: string;
	
	    static createFrom(source: any = {}) {
	        return new SourceOutline(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.params = source["params"];
	        this.returns = source["returns"];
	        this.loops = this.convertValues(source["loops"], OutlinePoint);
	        this.branches = this.convertValues(source["branches"], OutlinePoint);
	        this.returnsAt = this.convertValues(source["returnsAt"], OutlinePoint);
	        this.demoCall = source["demoCall"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class AnalyzeResult {
	    outline: SourceOutline;
	    summary: string;
	
	    static createFrom(source: any = {}) {
	        return new AnalyzeResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.outline = this.convertValues(source["outline"], SourceOutline);
	        this.summary = source["summary"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class AppSettings {
	    language: string;
	    theme: string;
	    uiStyle: string;
	    accentColor: string;
	    customCodeDir: string;
	    dataDir: string;
	    playbackSpeed: number;
	
	    static createFrom(source: any = {}) {
	        return new AppSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.language = source["language"];
	        this.theme = source["theme"];
	        this.uiStyle = source["uiStyle"];
	        this.accentColor = source["accentColor"];
	        this.customCodeDir = source["customCodeDir"];
	        this.dataDir = source["dataDir"];
	        this.playbackSpeed = source["playbackSpeed"];
	    }
	}
	export class CompileError {
	    file: string;
	    line: number;
	    column: number;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new CompileError(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.file = source["file"];
	        this.line = source["line"];
	        this.column = source["column"];
	        this.message = source["message"];
	    }
	}
	export class CompileResult {
	    success: boolean;
	    sessionId: string;
	    binaryPath?: string;
	    compiledSource?: string;
	    debugSource?: string;
	    debugFile?: string;
	    editorLineOffset: number;
	    errors?: CompileError[];
	
	    static createFrom(source: any = {}) {
	        return new CompileResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.sessionId = source["sessionId"];
	        this.binaryPath = source["binaryPath"];
	        this.compiledSource = source["compiledSource"];
	        this.debugSource = source["debugSource"];
	        this.debugFile = source["debugFile"];
	        this.editorLineOffset = source["editorLineOffset"];
	        this.errors = this.convertValues(source["errors"], CompileError);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class VarDelta {
	    name: string;
	    type: string;
	    oldValue: string;
	    newValue: string;
	
	    static createFrom(source: any = {}) {
	        return new VarDelta(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.type = source["type"];
	        this.oldValue = source["oldValue"];
	        this.newValue = source["newValue"];
	    }
	}
	export class FrameExplanation {
	    summary: string;
	    lineText: string;
	    stmtKind: string;
	    condition?: string;
	    changes: VarDelta[];
	    loopIter?: number;
	
	    static createFrom(source: any = {}) {
	        return new FrameExplanation(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.summary = source["summary"];
	        this.lineText = source["lineText"];
	        this.stmtKind = source["stmtKind"];
	        this.condition = source["condition"];
	        this.changes = this.convertValues(source["changes"], VarDelta);
	        this.loopIter = source["loopIter"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class TestCase {
	    name?: string;
	    input: string;
	    args?: string;
	    setup?: string;
	    expected: string;
	    hidden?: boolean;
	
	    static createFrom(source: any = {}) {
	        return new TestCase(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.input = source["input"];
	        this.args = source["args"];
	        this.setup = source["setup"];
	        this.expected = source["expected"];
	        this.hidden = source["hidden"];
	    }
	}
	export class Problem {
	    id: string;
	    title: string;
	    tags: string[];
	    template: string;
	    entry: string;
	    runner: string;
	    tests: TestCase[];
	
	    static createFrom(source: any = {}) {
	        return new Problem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.tags = source["tags"];
	        this.template = source["template"];
	        this.entry = source["entry"];
	        this.runner = source["runner"];
	        this.tests = this.convertValues(source["tests"], TestCase);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ProblemMeta {
	    id: string;
	    title: string;
	    tags: string[];
	
	    static createFrom(source: any = {}) {
	        return new ProblemMeta(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.tags = source["tags"];
	    }
	}
	export class RunOptions {
	    maxSteps: number;
	    lang?: string;
	
	    static createFrom(source: any = {}) {
	        return new RunOptions(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.maxSteps = source["maxSteps"];
	        this.lang = source["lang"];
	    }
	}
	
	export class StackFrame {
	    file: string;
	    line: number;
	    function: string;
	
	    static createFrom(source: any = {}) {
	        return new StackFrame(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.file = source["file"];
	        this.line = source["line"];
	        this.function = source["function"];
	    }
	}
	export class Trace {
	    sessionId: string;
	    frames: TraceFrame[];
	    status: string;
	    panicMsg?: string;
	
	    static createFrom(source: any = {}) {
	        return new Trace(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.sessionId = source["sessionId"];
	        this.frames = this.convertValues(source["frames"], TraceFrame);
	        this.status = source["status"];
	        this.panicMsg = source["panicMsg"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class VizEntry {
	    key: string;
	    value: string;
	
	    static createFrom(source: any = {}) {
	        return new VizEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.key = source["key"];
	        this.value = source["value"];
	    }
	}
	export class VizEdge {
	    from: string;
	    to: string;
	
	    static createFrom(source: any = {}) {
	        return new VizEdge(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.from = source["from"];
	        this.to = source["to"];
	    }
	}
	export class VizNode {
	    id: string;
	    label: string;
	    x?: number;
	    y?: number;
	    active?: boolean;
	
	    static createFrom(source: any = {}) {
	        return new VizNode(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.label = source["label"];
	        this.x = source["x"];
	        this.y = source["y"];
	        this.active = source["active"];
	    }
	}
	export class VizPointer {
	    name: string;
	    index: number;
	
	    static createFrom(source: any = {}) {
	        return new VizPointer(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.index = source["index"];
	    }
	}
	export class VizRange {
	    from: number;
	    to: number;
	
	    static createFrom(source: any = {}) {
	        return new VizRange(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.from = source["from"];
	        this.to = source["to"];
	    }
	}
	export class VizModel {
	    type: string;
	    name: string;
	    values?: string[];
	    rows?: string[][];
	    highlight?: number[];
	    highlightNode?: string;
	    rangeHighlight?: VizRange;
	    pointers?: VizPointer[];
	    nodes?: VizNode[];
	    edges?: VizEdge[];
	    entries?: VizEntry[];
	
	    static createFrom(source: any = {}) {
	        return new VizModel(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.name = source["name"];
	        this.values = source["values"];
	        this.rows = source["rows"];
	        this.highlight = source["highlight"];
	        this.highlightNode = source["highlightNode"];
	        this.rangeHighlight = this.convertValues(source["rangeHighlight"], VizRange);
	        this.pointers = this.convertValues(source["pointers"], VizPointer);
	        this.nodes = this.convertValues(source["nodes"], VizNode);
	        this.edges = this.convertValues(source["edges"], VizEdge);
	        this.entries = this.convertValues(source["entries"], VizEntry);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Variable {
	    name: string;
	    type: string;
	    value: string;
	    children?: Variable[];
	
	    static createFrom(source: any = {}) {
	        return new Variable(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.type = source["type"];
	        this.value = source["value"];
	        this.children = this.convertValues(source["children"], Variable);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class TraceFrame {
	    index: number;
	    file: string;
	    line: number;
	    editorLine: number;
	    function: string;
	    variables: Variable[];
	    stack: StackFrame[];
	    viz: VizModel[];
	    explanation?: FrameExplanation;
	
	    static createFrom(source: any = {}) {
	        return new TraceFrame(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.index = source["index"];
	        this.file = source["file"];
	        this.line = source["line"];
	        this.editorLine = source["editorLine"];
	        this.function = source["function"];
	        this.variables = this.convertValues(source["variables"], Variable);
	        this.stack = this.convertValues(source["stack"], StackFrame);
	        this.viz = this.convertValues(source["viz"], VizModel);
	        this.explanation = this.convertValues(source["explanation"], FrameExplanation);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class StepResult {
	    frame?: TraceFrame;
	    trace?: Trace;
	    finished: boolean;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new StepResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.frame = this.convertValues(source["frame"], TraceFrame);
	        this.trace = this.convertValues(source["trace"], Trace);
	        this.finished = source["finished"];
	        this.error = source["error"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class TestResult {
	    index: number;
	    passed: boolean;
	    actual?: string;
	    expected: string;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new TestResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.index = source["index"];
	        this.passed = source["passed"];
	        this.actual = source["actual"];
	        this.expected = source["expected"];
	        this.error = source["error"];
	    }
	}
	
	
	
	
	
	
	
	
	

}

