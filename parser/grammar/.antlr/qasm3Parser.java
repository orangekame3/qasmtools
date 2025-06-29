// Generated from /Users/orangekame3/src/github.com/orangekame3/qasmtools/parser/grammar/qasm3Parser.g4 by ANTLR 4.13.1
import org.antlr.v4.runtime.atn.*;
import org.antlr.v4.runtime.dfa.DFA;
import org.antlr.v4.runtime.*;
import org.antlr.v4.runtime.misc.*;
import org.antlr.v4.runtime.tree.*;
import java.util.List;
import java.util.Iterator;
import java.util.ArrayList;

@SuppressWarnings({"all", "warnings", "unchecked", "unused", "cast", "CheckReturnValue"})
public class qasm3Parser extends Parser {
	static { RuntimeMetaData.checkVersion("4.13.1", RuntimeMetaData.VERSION); }

	protected static final DFA[] _decisionToDFA;
	protected static final PredictionContextCache _sharedContextCache =
		new PredictionContextCache();
	public static final int
		OPENQASM=1, INCLUDE=2, DEFCALGRAMMAR=3, DEF=4, CAL=5, DEFCAL=6, GATE=7, 
		EXTERN=8, BOX=9, LET=10, BREAK=11, CONTINUE=12, IF=13, ELSE=14, END=15, 
		RETURN=16, FOR=17, WHILE=18, IN=19, SWITCH=20, CASE=21, DEFAULT=22, NOP=23, 
		PRAGMA=24, AnnotationKeyword=25, INPUT=26, OUTPUT=27, CONST=28, READONLY=29, 
		MUTABLE=30, QREG=31, QUBIT=32, CREG=33, BOOL=34, BIT=35, INT=36, UINT=37, 
		FLOAT=38, ANGLE=39, COMPLEX=40, ARRAY=41, VOID=42, DURATION=43, STRETCH=44, 
		GPHASE=45, INV=46, POW=47, CTRL=48, NEGCTRL=49, DIM=50, DURATIONOF=51, 
		DELAY=52, RESET=53, MEASURE=54, BARRIER=55, BooleanLiteral=56, LBRACKET=57, 
		RBRACKET=58, LBRACE=59, RBRACE=60, LPAREN=61, RPAREN=62, COLON=63, SEMICOLON=64, 
		DOT=65, COMMA=66, EQUALS=67, ARROW=68, PLUS=69, DOUBLE_PLUS=70, MINUS=71, 
		ASTERISK=72, DOUBLE_ASTERISK=73, SLASH=74, PERCENT=75, PIPE=76, DOUBLE_PIPE=77, 
		AMPERSAND=78, DOUBLE_AMPERSAND=79, CARET=80, AT=81, TILDE=82, EXCLAMATION_POINT=83, 
		EqualityOperator=84, CompoundAssignmentOperator=85, ComparisonOperator=86, 
		BitshiftOperator=87, IMAG=88, ImaginaryLiteral=89, BinaryIntegerLiteral=90, 
		OctalIntegerLiteral=91, DecimalIntegerLiteral=92, HexIntegerLiteral=93, 
		Identifier=94, HardwareQubit=95, FloatLiteral=96, TimingLiteral=97, BitstringLiteral=98, 
		Whitespace=99, Newline=100, LineComment=101, BlockComment=102, VERSION_IDENTIFER_WHITESPACE=103, 
		VersionSpecifier=104, ARBITRARY_STRING_WHITESPACE=105, StringLiteral=106, 
		EAT_INITIAL_SPACE=107, EAT_LINE_END=108, RemainingLineContent=109, CAL_PRELUDE_WHITESPACE=110, 
		CAL_PRELUDE_COMMENT=111, DEFCAL_PRELUDE_WHITESPACE=112, DEFCAL_PRELUDE_COMMENT=113, 
		CalibrationBlock=114;
	public static final int
		RULE_program = 0, RULE_version = 1, RULE_statement = 2, RULE_annotation = 3, 
		RULE_scope = 4, RULE_pragma = 5, RULE_statementOrScope = 6, RULE_calibrationGrammarStatement = 7, 
		RULE_includeStatement = 8, RULE_breakStatement = 9, RULE_continueStatement = 10, 
		RULE_endStatement = 11, RULE_forStatement = 12, RULE_ifStatement = 13, 
		RULE_returnStatement = 14, RULE_whileStatement = 15, RULE_switchStatement = 16, 
		RULE_switchCaseItem = 17, RULE_barrierStatement = 18, RULE_boxStatement = 19, 
		RULE_delayStatement = 20, RULE_nopStatement = 21, RULE_gateCallStatement = 22, 
		RULE_measureArrowAssignmentStatement = 23, RULE_resetStatement = 24, RULE_aliasDeclarationStatement = 25, 
		RULE_classicalDeclarationStatement = 26, RULE_constDeclarationStatement = 27, 
		RULE_ioDeclarationStatement = 28, RULE_oldStyleDeclarationStatement = 29, 
		RULE_quantumDeclarationStatement = 30, RULE_defStatement = 31, RULE_externStatement = 32, 
		RULE_gateStatement = 33, RULE_assignmentStatement = 34, RULE_expressionStatement = 35, 
		RULE_calStatement = 36, RULE_defcalStatement = 37, RULE_expression = 38, 
		RULE_aliasExpression = 39, RULE_declarationExpression = 40, RULE_measureExpression = 41, 
		RULE_rangeExpression = 42, RULE_setExpression = 43, RULE_arrayLiteral = 44, 
		RULE_indexOperator = 45, RULE_indexedIdentifier = 46, RULE_returnSignature = 47, 
		RULE_gateModifier = 48, RULE_scalarType = 49, RULE_qubitType = 50, RULE_arrayType = 51, 
		RULE_arrayReferenceType = 52, RULE_designator = 53, RULE_defcalTarget = 54, 
		RULE_defcalArgumentDefinition = 55, RULE_defcalOperand = 56, RULE_gateOperand = 57, 
		RULE_externArgument = 58, RULE_argumentDefinition = 59, RULE_argumentDefinitionList = 60, 
		RULE_defcalArgumentDefinitionList = 61, RULE_defcalOperandList = 62, RULE_expressionList = 63, 
		RULE_identifierList = 64, RULE_gateOperandList = 65, RULE_externArgumentList = 66;
	private static String[] makeRuleNames() {
		return new String[] {
			"program", "version", "statement", "annotation", "scope", "pragma", "statementOrScope", 
			"calibrationGrammarStatement", "includeStatement", "breakStatement", 
			"continueStatement", "endStatement", "forStatement", "ifStatement", "returnStatement", 
			"whileStatement", "switchStatement", "switchCaseItem", "barrierStatement", 
			"boxStatement", "delayStatement", "nopStatement", "gateCallStatement", 
			"measureArrowAssignmentStatement", "resetStatement", "aliasDeclarationStatement", 
			"classicalDeclarationStatement", "constDeclarationStatement", "ioDeclarationStatement", 
			"oldStyleDeclarationStatement", "quantumDeclarationStatement", "defStatement", 
			"externStatement", "gateStatement", "assignmentStatement", "expressionStatement", 
			"calStatement", "defcalStatement", "expression", "aliasExpression", "declarationExpression", 
			"measureExpression", "rangeExpression", "setExpression", "arrayLiteral", 
			"indexOperator", "indexedIdentifier", "returnSignature", "gateModifier", 
			"scalarType", "qubitType", "arrayType", "arrayReferenceType", "designator", 
			"defcalTarget", "defcalArgumentDefinition", "defcalOperand", "gateOperand", 
			"externArgument", "argumentDefinition", "argumentDefinitionList", "defcalArgumentDefinitionList", 
			"defcalOperandList", "expressionList", "identifierList", "gateOperandList", 
			"externArgumentList"
		};
	}
	public static final String[] ruleNames = makeRuleNames();

	private static String[] makeLiteralNames() {
		return new String[] {
			null, "'OPENQASM'", "'include'", "'defcalgrammar'", "'def'", "'cal'", 
			"'defcal'", "'gate'", "'extern'", "'box'", "'let'", "'break'", "'continue'", 
			"'if'", "'else'", "'end'", "'return'", "'for'", "'while'", "'in'", "'switch'", 
			"'case'", "'default'", "'nop'", null, null, "'input'", "'output'", "'const'", 
			"'readonly'", "'mutable'", "'qreg'", "'qubit'", "'creg'", "'bool'", "'bit'", 
			"'int'", "'uint'", "'float'", "'angle'", "'complex'", "'array'", "'void'", 
			"'duration'", "'stretch'", "'gphase'", "'inv'", "'pow'", "'ctrl'", "'negctrl'", 
			"'#dim'", "'durationof'", "'delay'", "'reset'", "'measure'", "'barrier'", 
			null, "'['", "']'", "'{'", "'}'", "'('", "')'", "':'", "';'", "'.'", 
			"','", "'='", "'->'", "'+'", "'++'", "'-'", "'*'", "'**'", "'/'", "'%'", 
			"'|'", "'||'", "'&'", "'&&'", "'^'", "'@'", "'~'", "'!'", null, null, 
			null, null, "'im'"
		};
	}
	private static final String[] _LITERAL_NAMES = makeLiteralNames();
	private static String[] makeSymbolicNames() {
		return new String[] {
			null, "OPENQASM", "INCLUDE", "DEFCALGRAMMAR", "DEF", "CAL", "DEFCAL", 
			"GATE", "EXTERN", "BOX", "LET", "BREAK", "CONTINUE", "IF", "ELSE", "END", 
			"RETURN", "FOR", "WHILE", "IN", "SWITCH", "CASE", "DEFAULT", "NOP", "PRAGMA", 
			"AnnotationKeyword", "INPUT", "OUTPUT", "CONST", "READONLY", "MUTABLE", 
			"QREG", "QUBIT", "CREG", "BOOL", "BIT", "INT", "UINT", "FLOAT", "ANGLE", 
			"COMPLEX", "ARRAY", "VOID", "DURATION", "STRETCH", "GPHASE", "INV", "POW", 
			"CTRL", "NEGCTRL", "DIM", "DURATIONOF", "DELAY", "RESET", "MEASURE", 
			"BARRIER", "BooleanLiteral", "LBRACKET", "RBRACKET", "LBRACE", "RBRACE", 
			"LPAREN", "RPAREN", "COLON", "SEMICOLON", "DOT", "COMMA", "EQUALS", "ARROW", 
			"PLUS", "DOUBLE_PLUS", "MINUS", "ASTERISK", "DOUBLE_ASTERISK", "SLASH", 
			"PERCENT", "PIPE", "DOUBLE_PIPE", "AMPERSAND", "DOUBLE_AMPERSAND", "CARET", 
			"AT", "TILDE", "EXCLAMATION_POINT", "EqualityOperator", "CompoundAssignmentOperator", 
			"ComparisonOperator", "BitshiftOperator", "IMAG", "ImaginaryLiteral", 
			"BinaryIntegerLiteral", "OctalIntegerLiteral", "DecimalIntegerLiteral", 
			"HexIntegerLiteral", "Identifier", "HardwareQubit", "FloatLiteral", "TimingLiteral", 
			"BitstringLiteral", "Whitespace", "Newline", "LineComment", "BlockComment", 
			"VERSION_IDENTIFER_WHITESPACE", "VersionSpecifier", "ARBITRARY_STRING_WHITESPACE", 
			"StringLiteral", "EAT_INITIAL_SPACE", "EAT_LINE_END", "RemainingLineContent", 
			"CAL_PRELUDE_WHITESPACE", "CAL_PRELUDE_COMMENT", "DEFCAL_PRELUDE_WHITESPACE", 
			"DEFCAL_PRELUDE_COMMENT", "CalibrationBlock"
		};
	}
	private static final String[] _SYMBOLIC_NAMES = makeSymbolicNames();
	public static final Vocabulary VOCABULARY = new VocabularyImpl(_LITERAL_NAMES, _SYMBOLIC_NAMES);

	/**
	 * @deprecated Use {@link #VOCABULARY} instead.
	 */
	@Deprecated
	public static final String[] tokenNames;
	static {
		tokenNames = new String[_SYMBOLIC_NAMES.length];
		for (int i = 0; i < tokenNames.length; i++) {
			tokenNames[i] = VOCABULARY.getLiteralName(i);
			if (tokenNames[i] == null) {
				tokenNames[i] = VOCABULARY.getSymbolicName(i);
			}

			if (tokenNames[i] == null) {
				tokenNames[i] = "<INVALID>";
			}
		}
	}

	@Override
	@Deprecated
	public String[] getTokenNames() {
		return tokenNames;
	}

	@Override

	public Vocabulary getVocabulary() {
		return VOCABULARY;
	}

	@Override
	public String getGrammarFileName() { return "qasm3Parser.g4"; }

	@Override
	public String[] getRuleNames() { return ruleNames; }

	@Override
	public String getSerializedATN() { return _serializedATN; }

	@Override
	public ATN getATN() { return _ATN; }

	public qasm3Parser(TokenStream input) {
		super(input);
		_interp = new ParserATNSimulator(this,_ATN,_decisionToDFA,_sharedContextCache);
	}

	@SuppressWarnings("CheckReturnValue")
	public static class ProgramContext extends ParserRuleContext {
		public TerminalNode EOF() { return getToken(qasm3Parser.EOF, 0); }
		public VersionContext version() {
			return getRuleContext(VersionContext.class,0);
		}
		public List<StatementOrScopeContext> statementOrScope() {
			return getRuleContexts(StatementOrScopeContext.class);
		}
		public StatementOrScopeContext statementOrScope(int i) {
			return getRuleContext(StatementOrScopeContext.class,i);
		}
		public ProgramContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_program; }
	}

	public final ProgramContext program() throws RecognitionException {
		ProgramContext _localctx = new ProgramContext(_ctx, getState());
		enterRule(_localctx, 0, RULE_program);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(135);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if (_la==OPENQASM) {
				{
				setState(134);
				version();
				}
			}

			setState(140);
			_errHandler.sync(this);
			_la = _input.LA(1);
			while ((((_la) & ~0x3f) == 0 && ((1L << _la) & 3025288650022174716L) != 0) || ((((_la - 71)) & ~0x3f) == 0 && ((1L << (_la - 71)) & 268179457L) != 0)) {
				{
				{
				setState(137);
				statementOrScope();
				}
				}
				setState(142);
				_errHandler.sync(this);
				_la = _input.LA(1);
			}
			setState(143);
			match(EOF);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class VersionContext extends ParserRuleContext {
		public TerminalNode OPENQASM() { return getToken(qasm3Parser.OPENQASM, 0); }
		public TerminalNode VersionSpecifier() { return getToken(qasm3Parser.VersionSpecifier, 0); }
		public TerminalNode SEMICOLON() { return getToken(qasm3Parser.SEMICOLON, 0); }
		public VersionContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_version; }
	}

	public final VersionContext version() throws RecognitionException {
		VersionContext _localctx = new VersionContext(_ctx, getState());
		enterRule(_localctx, 2, RULE_version);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(145);
			match(OPENQASM);
			setState(146);
			match(VersionSpecifier);
			setState(147);
			match(SEMICOLON);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class StatementContext extends ParserRuleContext {
		public PragmaContext pragma() {
			return getRuleContext(PragmaContext.class,0);
		}
		public AliasDeclarationStatementContext aliasDeclarationStatement() {
			return getRuleContext(AliasDeclarationStatementContext.class,0);
		}
		public AssignmentStatementContext assignmentStatement() {
			return getRuleContext(AssignmentStatementContext.class,0);
		}
		public BarrierStatementContext barrierStatement() {
			return getRuleContext(BarrierStatementContext.class,0);
		}
		public BoxStatementContext boxStatement() {
			return getRuleContext(BoxStatementContext.class,0);
		}
		public BreakStatementContext breakStatement() {
			return getRuleContext(BreakStatementContext.class,0);
		}
		public CalStatementContext calStatement() {
			return getRuleContext(CalStatementContext.class,0);
		}
		public CalibrationGrammarStatementContext calibrationGrammarStatement() {
			return getRuleContext(CalibrationGrammarStatementContext.class,0);
		}
		public ClassicalDeclarationStatementContext classicalDeclarationStatement() {
			return getRuleContext(ClassicalDeclarationStatementContext.class,0);
		}
		public ConstDeclarationStatementContext constDeclarationStatement() {
			return getRuleContext(ConstDeclarationStatementContext.class,0);
		}
		public ContinueStatementContext continueStatement() {
			return getRuleContext(ContinueStatementContext.class,0);
		}
		public DefStatementContext defStatement() {
			return getRuleContext(DefStatementContext.class,0);
		}
		public DefcalStatementContext defcalStatement() {
			return getRuleContext(DefcalStatementContext.class,0);
		}
		public DelayStatementContext delayStatement() {
			return getRuleContext(DelayStatementContext.class,0);
		}
		public EndStatementContext endStatement() {
			return getRuleContext(EndStatementContext.class,0);
		}
		public ExpressionStatementContext expressionStatement() {
			return getRuleContext(ExpressionStatementContext.class,0);
		}
		public ExternStatementContext externStatement() {
			return getRuleContext(ExternStatementContext.class,0);
		}
		public ForStatementContext forStatement() {
			return getRuleContext(ForStatementContext.class,0);
		}
		public GateCallStatementContext gateCallStatement() {
			return getRuleContext(GateCallStatementContext.class,0);
		}
		public GateStatementContext gateStatement() {
			return getRuleContext(GateStatementContext.class,0);
		}
		public IfStatementContext ifStatement() {
			return getRuleContext(IfStatementContext.class,0);
		}
		public IncludeStatementContext includeStatement() {
			return getRuleContext(IncludeStatementContext.class,0);
		}
		public IoDeclarationStatementContext ioDeclarationStatement() {
			return getRuleContext(IoDeclarationStatementContext.class,0);
		}
		public MeasureArrowAssignmentStatementContext measureArrowAssignmentStatement() {
			return getRuleContext(MeasureArrowAssignmentStatementContext.class,0);
		}
		public NopStatementContext nopStatement() {
			return getRuleContext(NopStatementContext.class,0);
		}
		public OldStyleDeclarationStatementContext oldStyleDeclarationStatement() {
			return getRuleContext(OldStyleDeclarationStatementContext.class,0);
		}
		public QuantumDeclarationStatementContext quantumDeclarationStatement() {
			return getRuleContext(QuantumDeclarationStatementContext.class,0);
		}
		public ResetStatementContext resetStatement() {
			return getRuleContext(ResetStatementContext.class,0);
		}
		public ReturnStatementContext returnStatement() {
			return getRuleContext(ReturnStatementContext.class,0);
		}
		public SwitchStatementContext switchStatement() {
			return getRuleContext(SwitchStatementContext.class,0);
		}
		public WhileStatementContext whileStatement() {
			return getRuleContext(WhileStatementContext.class,0);
		}
		public List<AnnotationContext> annotation() {
			return getRuleContexts(AnnotationContext.class);
		}
		public AnnotationContext annotation(int i) {
			return getRuleContext(AnnotationContext.class,i);
		}
		public StatementContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_statement; }
	}

	public final StatementContext statement() throws RecognitionException {
		StatementContext _localctx = new StatementContext(_ctx, getState());
		enterRule(_localctx, 4, RULE_statement);
		int _la;
		try {
			setState(188);
			_errHandler.sync(this);
			switch (_input.LA(1)) {
			case PRAGMA:
				enterOuterAlt(_localctx, 1);
				{
				setState(149);
				pragma();
				}
				break;
			case INCLUDE:
			case DEFCALGRAMMAR:
			case DEF:
			case CAL:
			case DEFCAL:
			case GATE:
			case EXTERN:
			case BOX:
			case LET:
			case BREAK:
			case CONTINUE:
			case IF:
			case END:
			case RETURN:
			case FOR:
			case WHILE:
			case SWITCH:
			case NOP:
			case AnnotationKeyword:
			case INPUT:
			case OUTPUT:
			case CONST:
			case QREG:
			case QUBIT:
			case CREG:
			case BOOL:
			case BIT:
			case INT:
			case UINT:
			case FLOAT:
			case ANGLE:
			case COMPLEX:
			case ARRAY:
			case DURATION:
			case STRETCH:
			case GPHASE:
			case INV:
			case POW:
			case CTRL:
			case NEGCTRL:
			case DURATIONOF:
			case DELAY:
			case RESET:
			case MEASURE:
			case BARRIER:
			case BooleanLiteral:
			case LPAREN:
			case MINUS:
			case TILDE:
			case EXCLAMATION_POINT:
			case ImaginaryLiteral:
			case BinaryIntegerLiteral:
			case OctalIntegerLiteral:
			case DecimalIntegerLiteral:
			case HexIntegerLiteral:
			case Identifier:
			case HardwareQubit:
			case FloatLiteral:
			case TimingLiteral:
			case BitstringLiteral:
				enterOuterAlt(_localctx, 2);
				{
				setState(153);
				_errHandler.sync(this);
				_la = _input.LA(1);
				while (_la==AnnotationKeyword) {
					{
					{
					setState(150);
					annotation();
					}
					}
					setState(155);
					_errHandler.sync(this);
					_la = _input.LA(1);
				}
				setState(186);
				_errHandler.sync(this);
				switch ( getInterpreter().adaptivePredict(_input,3,_ctx) ) {
				case 1:
					{
					setState(156);
					aliasDeclarationStatement();
					}
					break;
				case 2:
					{
					setState(157);
					assignmentStatement();
					}
					break;
				case 3:
					{
					setState(158);
					barrierStatement();
					}
					break;
				case 4:
					{
					setState(159);
					boxStatement();
					}
					break;
				case 5:
					{
					setState(160);
					breakStatement();
					}
					break;
				case 6:
					{
					setState(161);
					calStatement();
					}
					break;
				case 7:
					{
					setState(162);
					calibrationGrammarStatement();
					}
					break;
				case 8:
					{
					setState(163);
					classicalDeclarationStatement();
					}
					break;
				case 9:
					{
					setState(164);
					constDeclarationStatement();
					}
					break;
				case 10:
					{
					setState(165);
					continueStatement();
					}
					break;
				case 11:
					{
					setState(166);
					defStatement();
					}
					break;
				case 12:
					{
					setState(167);
					defcalStatement();
					}
					break;
				case 13:
					{
					setState(168);
					delayStatement();
					}
					break;
				case 14:
					{
					setState(169);
					endStatement();
					}
					break;
				case 15:
					{
					setState(170);
					expressionStatement();
					}
					break;
				case 16:
					{
					setState(171);
					externStatement();
					}
					break;
				case 17:
					{
					setState(172);
					forStatement();
					}
					break;
				case 18:
					{
					setState(173);
					gateCallStatement();
					}
					break;
				case 19:
					{
					setState(174);
					gateStatement();
					}
					break;
				case 20:
					{
					setState(175);
					ifStatement();
					}
					break;
				case 21:
					{
					setState(176);
					includeStatement();
					}
					break;
				case 22:
					{
					setState(177);
					ioDeclarationStatement();
					}
					break;
				case 23:
					{
					setState(178);
					measureArrowAssignmentStatement();
					}
					break;
				case 24:
					{
					setState(179);
					nopStatement();
					}
					break;
				case 25:
					{
					setState(180);
					oldStyleDeclarationStatement();
					}
					break;
				case 26:
					{
					setState(181);
					quantumDeclarationStatement();
					}
					break;
				case 27:
					{
					setState(182);
					resetStatement();
					}
					break;
				case 28:
					{
					setState(183);
					returnStatement();
					}
					break;
				case 29:
					{
					setState(184);
					switchStatement();
					}
					break;
				case 30:
					{
					setState(185);
					whileStatement();
					}
					break;
				}
				}
				break;
			default:
				throw new NoViableAltException(this);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class AnnotationContext extends ParserRuleContext {
		public TerminalNode AnnotationKeyword() { return getToken(qasm3Parser.AnnotationKeyword, 0); }
		public TerminalNode RemainingLineContent() { return getToken(qasm3Parser.RemainingLineContent, 0); }
		public AnnotationContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_annotation; }
	}

	public final AnnotationContext annotation() throws RecognitionException {
		AnnotationContext _localctx = new AnnotationContext(_ctx, getState());
		enterRule(_localctx, 6, RULE_annotation);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(190);
			match(AnnotationKeyword);
			setState(192);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if (_la==RemainingLineContent) {
				{
				setState(191);
				match(RemainingLineContent);
				}
			}

			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class ScopeContext extends ParserRuleContext {
		public TerminalNode LBRACE() { return getToken(qasm3Parser.LBRACE, 0); }
		public TerminalNode RBRACE() { return getToken(qasm3Parser.RBRACE, 0); }
		public List<StatementOrScopeContext> statementOrScope() {
			return getRuleContexts(StatementOrScopeContext.class);
		}
		public StatementOrScopeContext statementOrScope(int i) {
			return getRuleContext(StatementOrScopeContext.class,i);
		}
		public ScopeContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_scope; }
	}

	public final ScopeContext scope() throws RecognitionException {
		ScopeContext _localctx = new ScopeContext(_ctx, getState());
		enterRule(_localctx, 8, RULE_scope);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(194);
			match(LBRACE);
			setState(198);
			_errHandler.sync(this);
			_la = _input.LA(1);
			while ((((_la) & ~0x3f) == 0 && ((1L << _la) & 3025288650022174716L) != 0) || ((((_la - 71)) & ~0x3f) == 0 && ((1L << (_la - 71)) & 268179457L) != 0)) {
				{
				{
				setState(195);
				statementOrScope();
				}
				}
				setState(200);
				_errHandler.sync(this);
				_la = _input.LA(1);
			}
			setState(201);
			match(RBRACE);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class PragmaContext extends ParserRuleContext {
		public TerminalNode PRAGMA() { return getToken(qasm3Parser.PRAGMA, 0); }
		public TerminalNode RemainingLineContent() { return getToken(qasm3Parser.RemainingLineContent, 0); }
		public PragmaContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_pragma; }
	}

	public final PragmaContext pragma() throws RecognitionException {
		PragmaContext _localctx = new PragmaContext(_ctx, getState());
		enterRule(_localctx, 10, RULE_pragma);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(203);
			match(PRAGMA);
			setState(204);
			match(RemainingLineContent);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class StatementOrScopeContext extends ParserRuleContext {
		public StatementContext statement() {
			return getRuleContext(StatementContext.class,0);
		}
		public ScopeContext scope() {
			return getRuleContext(ScopeContext.class,0);
		}
		public StatementOrScopeContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_statementOrScope; }
	}

	public final StatementOrScopeContext statementOrScope() throws RecognitionException {
		StatementOrScopeContext _localctx = new StatementOrScopeContext(_ctx, getState());
		enterRule(_localctx, 12, RULE_statementOrScope);
		try {
			setState(208);
			_errHandler.sync(this);
			switch (_input.LA(1)) {
			case INCLUDE:
			case DEFCALGRAMMAR:
			case DEF:
			case CAL:
			case DEFCAL:
			case GATE:
			case EXTERN:
			case BOX:
			case LET:
			case BREAK:
			case CONTINUE:
			case IF:
			case END:
			case RETURN:
			case FOR:
			case WHILE:
			case SWITCH:
			case NOP:
			case PRAGMA:
			case AnnotationKeyword:
			case INPUT:
			case OUTPUT:
			case CONST:
			case QREG:
			case QUBIT:
			case CREG:
			case BOOL:
			case BIT:
			case INT:
			case UINT:
			case FLOAT:
			case ANGLE:
			case COMPLEX:
			case ARRAY:
			case DURATION:
			case STRETCH:
			case GPHASE:
			case INV:
			case POW:
			case CTRL:
			case NEGCTRL:
			case DURATIONOF:
			case DELAY:
			case RESET:
			case MEASURE:
			case BARRIER:
			case BooleanLiteral:
			case LPAREN:
			case MINUS:
			case TILDE:
			case EXCLAMATION_POINT:
			case ImaginaryLiteral:
			case BinaryIntegerLiteral:
			case OctalIntegerLiteral:
			case DecimalIntegerLiteral:
			case HexIntegerLiteral:
			case Identifier:
			case HardwareQubit:
			case FloatLiteral:
			case TimingLiteral:
			case BitstringLiteral:
				enterOuterAlt(_localctx, 1);
				{
				setState(206);
				statement();
				}
				break;
			case LBRACE:
				enterOuterAlt(_localctx, 2);
				{
				setState(207);
				scope();
				}
				break;
			default:
				throw new NoViableAltException(this);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class CalibrationGrammarStatementContext extends ParserRuleContext {
		public TerminalNode DEFCALGRAMMAR() { return getToken(qasm3Parser.DEFCALGRAMMAR, 0); }
		public TerminalNode StringLiteral() { return getToken(qasm3Parser.StringLiteral, 0); }
		public TerminalNode SEMICOLON() { return getToken(qasm3Parser.SEMICOLON, 0); }
		public CalibrationGrammarStatementContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_calibrationGrammarStatement; }
	}

	public final CalibrationGrammarStatementContext calibrationGrammarStatement() throws RecognitionException {
		CalibrationGrammarStatementContext _localctx = new CalibrationGrammarStatementContext(_ctx, getState());
		enterRule(_localctx, 14, RULE_calibrationGrammarStatement);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(210);
			match(DEFCALGRAMMAR);
			setState(211);
			match(StringLiteral);
			setState(212);
			match(SEMICOLON);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class IncludeStatementContext extends ParserRuleContext {
		public TerminalNode INCLUDE() { return getToken(qasm3Parser.INCLUDE, 0); }
		public TerminalNode StringLiteral() { return getToken(qasm3Parser.StringLiteral, 0); }
		public TerminalNode SEMICOLON() { return getToken(qasm3Parser.SEMICOLON, 0); }
		public IncludeStatementContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_includeStatement; }
	}

	public final IncludeStatementContext includeStatement() throws RecognitionException {
		IncludeStatementContext _localctx = new IncludeStatementContext(_ctx, getState());
		enterRule(_localctx, 16, RULE_includeStatement);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(214);
			match(INCLUDE);
			setState(215);
			match(StringLiteral);
			setState(216);
			match(SEMICOLON);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class BreakStatementContext extends ParserRuleContext {
		public TerminalNode BREAK() { return getToken(qasm3Parser.BREAK, 0); }
		public TerminalNode SEMICOLON() { return getToken(qasm3Parser.SEMICOLON, 0); }
		public BreakStatementContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_breakStatement; }
	}

	public final BreakStatementContext breakStatement() throws RecognitionException {
		BreakStatementContext _localctx = new BreakStatementContext(_ctx, getState());
		enterRule(_localctx, 18, RULE_breakStatement);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(218);
			match(BREAK);
			setState(219);
			match(SEMICOLON);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class ContinueStatementContext extends ParserRuleContext {
		public TerminalNode CONTINUE() { return getToken(qasm3Parser.CONTINUE, 0); }
		public TerminalNode SEMICOLON() { return getToken(qasm3Parser.SEMICOLON, 0); }
		public ContinueStatementContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_continueStatement; }
	}

	public final ContinueStatementContext continueStatement() throws RecognitionException {
		ContinueStatementContext _localctx = new ContinueStatementContext(_ctx, getState());
		enterRule(_localctx, 20, RULE_continueStatement);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(221);
			match(CONTINUE);
			setState(222);
			match(SEMICOLON);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class EndStatementContext extends ParserRuleContext {
		public TerminalNode END() { return getToken(qasm3Parser.END, 0); }
		public TerminalNode SEMICOLON() { return getToken(qasm3Parser.SEMICOLON, 0); }
		public EndStatementContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_endStatement; }
	}

	public final EndStatementContext endStatement() throws RecognitionException {
		EndStatementContext _localctx = new EndStatementContext(_ctx, getState());
		enterRule(_localctx, 22, RULE_endStatement);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(224);
			match(END);
			setState(225);
			match(SEMICOLON);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class ForStatementContext extends ParserRuleContext {
		public StatementOrScopeContext body;
		public TerminalNode FOR() { return getToken(qasm3Parser.FOR, 0); }
		public ScalarTypeContext scalarType() {
			return getRuleContext(ScalarTypeContext.class,0);
		}
		public TerminalNode Identifier() { return getToken(qasm3Parser.Identifier, 0); }
		public TerminalNode IN() { return getToken(qasm3Parser.IN, 0); }
		public StatementOrScopeContext statementOrScope() {
			return getRuleContext(StatementOrScopeContext.class,0);
		}
		public SetExpressionContext setExpression() {
			return getRuleContext(SetExpressionContext.class,0);
		}
		public TerminalNode LBRACKET() { return getToken(qasm3Parser.LBRACKET, 0); }
		public RangeExpressionContext rangeExpression() {
			return getRuleContext(RangeExpressionContext.class,0);
		}
		public TerminalNode RBRACKET() { return getToken(qasm3Parser.RBRACKET, 0); }
		public ExpressionContext expression() {
			return getRuleContext(ExpressionContext.class,0);
		}
		public ForStatementContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_forStatement; }
	}

	public final ForStatementContext forStatement() throws RecognitionException {
		ForStatementContext _localctx = new ForStatementContext(_ctx, getState());
		enterRule(_localctx, 24, RULE_forStatement);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(227);
			match(FOR);
			setState(228);
			scalarType();
			setState(229);
			match(Identifier);
			setState(230);
			match(IN);
			setState(237);
			_errHandler.sync(this);
			switch (_input.LA(1)) {
			case LBRACE:
				{
				setState(231);
				setExpression();
				}
				break;
			case LBRACKET:
				{
				setState(232);
				match(LBRACKET);
				setState(233);
				rangeExpression();
				setState(234);
				match(RBRACKET);
				}
				break;
			case BOOL:
			case BIT:
			case INT:
			case UINT:
			case FLOAT:
			case ANGLE:
			case COMPLEX:
			case ARRAY:
			case DURATION:
			case STRETCH:
			case DURATIONOF:
			case BooleanLiteral:
			case LPAREN:
			case MINUS:
			case TILDE:
			case EXCLAMATION_POINT:
			case ImaginaryLiteral:
			case BinaryIntegerLiteral:
			case OctalIntegerLiteral:
			case DecimalIntegerLiteral:
			case HexIntegerLiteral:
			case Identifier:
			case HardwareQubit:
			case FloatLiteral:
			case TimingLiteral:
			case BitstringLiteral:
				{
				setState(236);
				expression(0);
				}
				break;
			default:
				throw new NoViableAltException(this);
			}
			setState(239);
			((ForStatementContext)_localctx).body = statementOrScope();
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class IfStatementContext extends ParserRuleContext {
		public StatementOrScopeContext if_body;
		public StatementOrScopeContext else_body;
		public TerminalNode IF() { return getToken(qasm3Parser.IF, 0); }
		public TerminalNode LPAREN() { return getToken(qasm3Parser.LPAREN, 0); }
		public ExpressionContext expression() {
			return getRuleContext(ExpressionContext.class,0);
		}
		public TerminalNode RPAREN() { return getToken(qasm3Parser.RPAREN, 0); }
		public List<StatementOrScopeContext> statementOrScope() {
			return getRuleContexts(StatementOrScopeContext.class);
		}
		public StatementOrScopeContext statementOrScope(int i) {
			return getRuleContext(StatementOrScopeContext.class,i);
		}
		public TerminalNode ELSE() { return getToken(qasm3Parser.ELSE, 0); }
		public IfStatementContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_ifStatement; }
	}

	public final IfStatementContext ifStatement() throws RecognitionException {
		IfStatementContext _localctx = new IfStatementContext(_ctx, getState());
		enterRule(_localctx, 26, RULE_ifStatement);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(241);
			match(IF);
			setState(242);
			match(LPAREN);
			setState(243);
			expression(0);
			setState(244);
			match(RPAREN);
			setState(245);
			((IfStatementContext)_localctx).if_body = statementOrScope();
			setState(248);
			_errHandler.sync(this);
			switch ( getInterpreter().adaptivePredict(_input,9,_ctx) ) {
			case 1:
				{
				setState(246);
				match(ELSE);
				setState(247);
				((IfStatementContext)_localctx).else_body = statementOrScope();
				}
				break;
			}
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class ReturnStatementContext extends ParserRuleContext {
		public TerminalNode RETURN() { return getToken(qasm3Parser.RETURN, 0); }
		public TerminalNode SEMICOLON() { return getToken(qasm3Parser.SEMICOLON, 0); }
		public ExpressionContext expression() {
			return getRuleContext(ExpressionContext.class,0);
		}
		public MeasureExpressionContext measureExpression() {
			return getRuleContext(MeasureExpressionContext.class,0);
		}
		public ReturnStatementContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_returnStatement; }
	}

	public final ReturnStatementContext returnStatement() throws RecognitionException {
		ReturnStatementContext _localctx = new ReturnStatementContext(_ctx, getState());
		enterRule(_localctx, 28, RULE_returnStatement);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(250);
			match(RETURN);
			setState(253);
			_errHandler.sync(this);
			switch (_input.LA(1)) {
			case BOOL:
			case BIT:
			case INT:
			case UINT:
			case FLOAT:
			case ANGLE:
			case COMPLEX:
			case ARRAY:
			case DURATION:
			case STRETCH:
			case DURATIONOF:
			case BooleanLiteral:
			case LPAREN:
			case MINUS:
			case TILDE:
			case EXCLAMATION_POINT:
			case ImaginaryLiteral:
			case BinaryIntegerLiteral:
			case OctalIntegerLiteral:
			case DecimalIntegerLiteral:
			case HexIntegerLiteral:
			case Identifier:
			case HardwareQubit:
			case FloatLiteral:
			case TimingLiteral:
			case BitstringLiteral:
				{
				setState(251);
				expression(0);
				}
				break;
			case MEASURE:
				{
				setState(252);
				measureExpression();
				}
				break;
			case SEMICOLON:
				break;
			default:
				break;
			}
			setState(255);
			match(SEMICOLON);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class WhileStatementContext extends ParserRuleContext {
		public StatementOrScopeContext body;
		public TerminalNode WHILE() { return getToken(qasm3Parser.WHILE, 0); }
		public TerminalNode LPAREN() { return getToken(qasm3Parser.LPAREN, 0); }
		public ExpressionContext expression() {
			return getRuleContext(ExpressionContext.class,0);
		}
		public TerminalNode RPAREN() { return getToken(qasm3Parser.RPAREN, 0); }
		public StatementOrScopeContext statementOrScope() {
			return getRuleContext(StatementOrScopeContext.class,0);
		}
		public WhileStatementContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_whileStatement; }
	}

	public final WhileStatementContext whileStatement() throws RecognitionException {
		WhileStatementContext _localctx = new WhileStatementContext(_ctx, getState());
		enterRule(_localctx, 30, RULE_whileStatement);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(257);
			match(WHILE);
			setState(258);
			match(LPAREN);
			setState(259);
			expression(0);
			setState(260);
			match(RPAREN);
			setState(261);
			((WhileStatementContext)_localctx).body = statementOrScope();
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class SwitchStatementContext extends ParserRuleContext {
		public TerminalNode SWITCH() { return getToken(qasm3Parser.SWITCH, 0); }
		public TerminalNode LPAREN() { return getToken(qasm3Parser.LPAREN, 0); }
		public ExpressionContext expression() {
			return getRuleContext(ExpressionContext.class,0);
		}
		public TerminalNode RPAREN() { return getToken(qasm3Parser.RPAREN, 0); }
		public TerminalNode LBRACE() { return getToken(qasm3Parser.LBRACE, 0); }
		public TerminalNode RBRACE() { return getToken(qasm3Parser.RBRACE, 0); }
		public List<SwitchCaseItemContext> switchCaseItem() {
			return getRuleContexts(SwitchCaseItemContext.class);
		}
		public SwitchCaseItemContext switchCaseItem(int i) {
			return getRuleContext(SwitchCaseItemContext.class,i);
		}
		public SwitchStatementContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_switchStatement; }
	}

	public final SwitchStatementContext switchStatement() throws RecognitionException {
		SwitchStatementContext _localctx = new SwitchStatementContext(_ctx, getState());
		enterRule(_localctx, 32, RULE_switchStatement);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(263);
			match(SWITCH);
			setState(264);
			match(LPAREN);
			setState(265);
			expression(0);
			setState(266);
			match(RPAREN);
			setState(267);
			match(LBRACE);
			setState(271);
			_errHandler.sync(this);
			_la = _input.LA(1);
			while (_la==CASE || _la==DEFAULT) {
				{
				{
				setState(268);
				switchCaseItem();
				}
				}
				setState(273);
				_errHandler.sync(this);
				_la = _input.LA(1);
			}
			setState(274);
			match(RBRACE);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class SwitchCaseItemContext extends ParserRuleContext {
		public TerminalNode CASE() { return getToken(qasm3Parser.CASE, 0); }
		public ExpressionListContext expressionList() {
			return getRuleContext(ExpressionListContext.class,0);
		}
		public ScopeContext scope() {
			return getRuleContext(ScopeContext.class,0);
		}
		public TerminalNode DEFAULT() { return getToken(qasm3Parser.DEFAULT, 0); }
		public SwitchCaseItemContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_switchCaseItem; }
	}

	public final SwitchCaseItemContext switchCaseItem() throws RecognitionException {
		SwitchCaseItemContext _localctx = new SwitchCaseItemContext(_ctx, getState());
		enterRule(_localctx, 34, RULE_switchCaseItem);
		try {
			setState(282);
			_errHandler.sync(this);
			switch (_input.LA(1)) {
			case CASE:
				enterOuterAlt(_localctx, 1);
				{
				setState(276);
				match(CASE);
				setState(277);
				expressionList();
				setState(278);
				scope();
				}
				break;
			case DEFAULT:
				enterOuterAlt(_localctx, 2);
				{
				setState(280);
				match(DEFAULT);
				setState(281);
				scope();
				}
				break;
			default:
				throw new NoViableAltException(this);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class BarrierStatementContext extends ParserRuleContext {
		public TerminalNode BARRIER() { return getToken(qasm3Parser.BARRIER, 0); }
		public TerminalNode SEMICOLON() { return getToken(qasm3Parser.SEMICOLON, 0); }
		public GateOperandListContext gateOperandList() {
			return getRuleContext(GateOperandListContext.class,0);
		}
		public BarrierStatementContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_barrierStatement; }
	}

	public final BarrierStatementContext barrierStatement() throws RecognitionException {
		BarrierStatementContext _localctx = new BarrierStatementContext(_ctx, getState());
		enterRule(_localctx, 36, RULE_barrierStatement);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(284);
			match(BARRIER);
			setState(286);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if (_la==Identifier || _la==HardwareQubit) {
				{
				setState(285);
				gateOperandList();
				}
			}

			setState(288);
			match(SEMICOLON);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class BoxStatementContext extends ParserRuleContext {
		public TerminalNode BOX() { return getToken(qasm3Parser.BOX, 0); }
		public ScopeContext scope() {
			return getRuleContext(ScopeContext.class,0);
		}
		public DesignatorContext designator() {
			return getRuleContext(DesignatorContext.class,0);
		}
		public BoxStatementContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_boxStatement; }
	}

	public final BoxStatementContext boxStatement() throws RecognitionException {
		BoxStatementContext _localctx = new BoxStatementContext(_ctx, getState());
		enterRule(_localctx, 38, RULE_boxStatement);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(290);
			match(BOX);
			setState(292);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if (_la==LBRACKET) {
				{
				setState(291);
				designator();
				}
			}

			setState(294);
			scope();
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class DelayStatementContext extends ParserRuleContext {
		public TerminalNode DELAY() { return getToken(qasm3Parser.DELAY, 0); }
		public DesignatorContext designator() {
			return getRuleContext(DesignatorContext.class,0);
		}
		public TerminalNode SEMICOLON() { return getToken(qasm3Parser.SEMICOLON, 0); }
		public GateOperandListContext gateOperandList() {
			return getRuleContext(GateOperandListContext.class,0);
		}
		public DelayStatementContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_delayStatement; }
	}

	public final DelayStatementContext delayStatement() throws RecognitionException {
		DelayStatementContext _localctx = new DelayStatementContext(_ctx, getState());
		enterRule(_localctx, 40, RULE_delayStatement);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(296);
			match(DELAY);
			setState(297);
			designator();
			setState(299);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if (_la==Identifier || _la==HardwareQubit) {
				{
				setState(298);
				gateOperandList();
				}
			}

			setState(301);
			match(SEMICOLON);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class NopStatementContext extends ParserRuleContext {
		public TerminalNode NOP() { return getToken(qasm3Parser.NOP, 0); }
		public TerminalNode SEMICOLON() { return getToken(qasm3Parser.SEMICOLON, 0); }
		public GateOperandListContext gateOperandList() {
			return getRuleContext(GateOperandListContext.class,0);
		}
		public NopStatementContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_nopStatement; }
	}

	public final NopStatementContext nopStatement() throws RecognitionException {
		NopStatementContext _localctx = new NopStatementContext(_ctx, getState());
		enterRule(_localctx, 42, RULE_nopStatement);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(303);
			match(NOP);
			setState(305);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if (_la==Identifier || _la==HardwareQubit) {
				{
				setState(304);
				gateOperandList();
				}
			}

			setState(307);
			match(SEMICOLON);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class GateCallStatementContext extends ParserRuleContext {
		public TerminalNode Identifier() { return getToken(qasm3Parser.Identifier, 0); }
		public GateOperandListContext gateOperandList() {
			return getRuleContext(GateOperandListContext.class,0);
		}
		public TerminalNode SEMICOLON() { return getToken(qasm3Parser.SEMICOLON, 0); }
		public List<GateModifierContext> gateModifier() {
			return getRuleContexts(GateModifierContext.class);
		}
		public GateModifierContext gateModifier(int i) {
			return getRuleContext(GateModifierContext.class,i);
		}
		public TerminalNode LPAREN() { return getToken(qasm3Parser.LPAREN, 0); }
		public TerminalNode RPAREN() { return getToken(qasm3Parser.RPAREN, 0); }
		public DesignatorContext designator() {
			return getRuleContext(DesignatorContext.class,0);
		}
		public ExpressionListContext expressionList() {
			return getRuleContext(ExpressionListContext.class,0);
		}
		public TerminalNode GPHASE() { return getToken(qasm3Parser.GPHASE, 0); }
		public GateCallStatementContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_gateCallStatement; }
	}

	public final GateCallStatementContext gateCallStatement() throws RecognitionException {
		GateCallStatementContext _localctx = new GateCallStatementContext(_ctx, getState());
		enterRule(_localctx, 44, RULE_gateCallStatement);
		int _la;
		try {
			setState(350);
			_errHandler.sync(this);
			switch ( getInterpreter().adaptivePredict(_input,26,_ctx) ) {
			case 1:
				enterOuterAlt(_localctx, 1);
				{
				setState(312);
				_errHandler.sync(this);
				_la = _input.LA(1);
				while ((((_la) & ~0x3f) == 0 && ((1L << _la) & 1055531162664960L) != 0)) {
					{
					{
					setState(309);
					gateModifier();
					}
					}
					setState(314);
					_errHandler.sync(this);
					_la = _input.LA(1);
				}
				setState(315);
				match(Identifier);
				setState(321);
				_errHandler.sync(this);
				_la = _input.LA(1);
				if (_la==LPAREN) {
					{
					setState(316);
					match(LPAREN);
					setState(318);
					_errHandler.sync(this);
					_la = _input.LA(1);
					if ((((_la) & ~0x3f) == 0 && ((1L << _la) & 2380183172211015680L) != 0) || ((((_la - 71)) & ~0x3f) == 0 && ((1L << (_la - 71)) & 268179457L) != 0)) {
						{
						setState(317);
						expressionList();
						}
					}

					setState(320);
					match(RPAREN);
					}
				}

				setState(324);
				_errHandler.sync(this);
				_la = _input.LA(1);
				if (_la==LBRACKET) {
					{
					setState(323);
					designator();
					}
				}

				setState(326);
				gateOperandList();
				setState(327);
				match(SEMICOLON);
				}
				break;
			case 2:
				enterOuterAlt(_localctx, 2);
				{
				setState(332);
				_errHandler.sync(this);
				_la = _input.LA(1);
				while ((((_la) & ~0x3f) == 0 && ((1L << _la) & 1055531162664960L) != 0)) {
					{
					{
					setState(329);
					gateModifier();
					}
					}
					setState(334);
					_errHandler.sync(this);
					_la = _input.LA(1);
				}
				setState(335);
				match(GPHASE);
				setState(341);
				_errHandler.sync(this);
				_la = _input.LA(1);
				if (_la==LPAREN) {
					{
					setState(336);
					match(LPAREN);
					setState(338);
					_errHandler.sync(this);
					_la = _input.LA(1);
					if ((((_la) & ~0x3f) == 0 && ((1L << _la) & 2380183172211015680L) != 0) || ((((_la - 71)) & ~0x3f) == 0 && ((1L << (_la - 71)) & 268179457L) != 0)) {
						{
						setState(337);
						expressionList();
						}
					}

					setState(340);
					match(RPAREN);
					}
				}

				setState(344);
				_errHandler.sync(this);
				_la = _input.LA(1);
				if (_la==LBRACKET) {
					{
					setState(343);
					designator();
					}
				}

				setState(347);
				_errHandler.sync(this);
				_la = _input.LA(1);
				if (_la==Identifier || _la==HardwareQubit) {
					{
					setState(346);
					gateOperandList();
					}
				}

				setState(349);
				match(SEMICOLON);
				}
				break;
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class MeasureArrowAssignmentStatementContext extends ParserRuleContext {
		public MeasureExpressionContext measureExpression() {
			return getRuleContext(MeasureExpressionContext.class,0);
		}
		public TerminalNode SEMICOLON() { return getToken(qasm3Parser.SEMICOLON, 0); }
		public TerminalNode ARROW() { return getToken(qasm3Parser.ARROW, 0); }
		public IndexedIdentifierContext indexedIdentifier() {
			return getRuleContext(IndexedIdentifierContext.class,0);
		}
		public MeasureArrowAssignmentStatementContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_measureArrowAssignmentStatement; }
	}

	public final MeasureArrowAssignmentStatementContext measureArrowAssignmentStatement() throws RecognitionException {
		MeasureArrowAssignmentStatementContext _localctx = new MeasureArrowAssignmentStatementContext(_ctx, getState());
		enterRule(_localctx, 46, RULE_measureArrowAssignmentStatement);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(352);
			measureExpression();
			setState(355);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if (_la==ARROW) {
				{
				setState(353);
				match(ARROW);
				setState(354);
				indexedIdentifier();
				}
			}

			setState(357);
			match(SEMICOLON);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class ResetStatementContext extends ParserRuleContext {
		public TerminalNode RESET() { return getToken(qasm3Parser.RESET, 0); }
		public GateOperandContext gateOperand() {
			return getRuleContext(GateOperandContext.class,0);
		}
		public TerminalNode SEMICOLON() { return getToken(qasm3Parser.SEMICOLON, 0); }
		public ResetStatementContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_resetStatement; }
	}

	public final ResetStatementContext resetStatement() throws RecognitionException {
		ResetStatementContext _localctx = new ResetStatementContext(_ctx, getState());
		enterRule(_localctx, 48, RULE_resetStatement);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(359);
			match(RESET);
			setState(360);
			gateOperand();
			setState(361);
			match(SEMICOLON);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class AliasDeclarationStatementContext extends ParserRuleContext {
		public TerminalNode LET() { return getToken(qasm3Parser.LET, 0); }
		public TerminalNode Identifier() { return getToken(qasm3Parser.Identifier, 0); }
		public TerminalNode EQUALS() { return getToken(qasm3Parser.EQUALS, 0); }
		public AliasExpressionContext aliasExpression() {
			return getRuleContext(AliasExpressionContext.class,0);
		}
		public TerminalNode SEMICOLON() { return getToken(qasm3Parser.SEMICOLON, 0); }
		public AliasDeclarationStatementContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_aliasDeclarationStatement; }
	}

	public final AliasDeclarationStatementContext aliasDeclarationStatement() throws RecognitionException {
		AliasDeclarationStatementContext _localctx = new AliasDeclarationStatementContext(_ctx, getState());
		enterRule(_localctx, 50, RULE_aliasDeclarationStatement);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(363);
			match(LET);
			setState(364);
			match(Identifier);
			setState(365);
			match(EQUALS);
			setState(366);
			aliasExpression();
			setState(367);
			match(SEMICOLON);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class ClassicalDeclarationStatementContext extends ParserRuleContext {
		public TerminalNode Identifier() { return getToken(qasm3Parser.Identifier, 0); }
		public TerminalNode SEMICOLON() { return getToken(qasm3Parser.SEMICOLON, 0); }
		public ScalarTypeContext scalarType() {
			return getRuleContext(ScalarTypeContext.class,0);
		}
		public ArrayTypeContext arrayType() {
			return getRuleContext(ArrayTypeContext.class,0);
		}
		public TerminalNode EQUALS() { return getToken(qasm3Parser.EQUALS, 0); }
		public DeclarationExpressionContext declarationExpression() {
			return getRuleContext(DeclarationExpressionContext.class,0);
		}
		public ClassicalDeclarationStatementContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_classicalDeclarationStatement; }
	}

	public final ClassicalDeclarationStatementContext classicalDeclarationStatement() throws RecognitionException {
		ClassicalDeclarationStatementContext _localctx = new ClassicalDeclarationStatementContext(_ctx, getState());
		enterRule(_localctx, 52, RULE_classicalDeclarationStatement);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(371);
			_errHandler.sync(this);
			switch (_input.LA(1)) {
			case BOOL:
			case BIT:
			case INT:
			case UINT:
			case FLOAT:
			case ANGLE:
			case COMPLEX:
			case DURATION:
			case STRETCH:
				{
				setState(369);
				scalarType();
				}
				break;
			case ARRAY:
				{
				setState(370);
				arrayType();
				}
				break;
			default:
				throw new NoViableAltException(this);
			}
			setState(373);
			match(Identifier);
			setState(376);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if (_la==EQUALS) {
				{
				setState(374);
				match(EQUALS);
				setState(375);
				declarationExpression();
				}
			}

			setState(378);
			match(SEMICOLON);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class ConstDeclarationStatementContext extends ParserRuleContext {
		public TerminalNode CONST() { return getToken(qasm3Parser.CONST, 0); }
		public ScalarTypeContext scalarType() {
			return getRuleContext(ScalarTypeContext.class,0);
		}
		public TerminalNode Identifier() { return getToken(qasm3Parser.Identifier, 0); }
		public TerminalNode EQUALS() { return getToken(qasm3Parser.EQUALS, 0); }
		public DeclarationExpressionContext declarationExpression() {
			return getRuleContext(DeclarationExpressionContext.class,0);
		}
		public TerminalNode SEMICOLON() { return getToken(qasm3Parser.SEMICOLON, 0); }
		public ConstDeclarationStatementContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_constDeclarationStatement; }
	}

	public final ConstDeclarationStatementContext constDeclarationStatement() throws RecognitionException {
		ConstDeclarationStatementContext _localctx = new ConstDeclarationStatementContext(_ctx, getState());
		enterRule(_localctx, 54, RULE_constDeclarationStatement);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(380);
			match(CONST);
			setState(381);
			scalarType();
			setState(382);
			match(Identifier);
			setState(383);
			match(EQUALS);
			setState(384);
			declarationExpression();
			setState(385);
			match(SEMICOLON);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class IoDeclarationStatementContext extends ParserRuleContext {
		public TerminalNode Identifier() { return getToken(qasm3Parser.Identifier, 0); }
		public TerminalNode SEMICOLON() { return getToken(qasm3Parser.SEMICOLON, 0); }
		public TerminalNode INPUT() { return getToken(qasm3Parser.INPUT, 0); }
		public TerminalNode OUTPUT() { return getToken(qasm3Parser.OUTPUT, 0); }
		public ScalarTypeContext scalarType() {
			return getRuleContext(ScalarTypeContext.class,0);
		}
		public ArrayTypeContext arrayType() {
			return getRuleContext(ArrayTypeContext.class,0);
		}
		public IoDeclarationStatementContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_ioDeclarationStatement; }
	}

	public final IoDeclarationStatementContext ioDeclarationStatement() throws RecognitionException {
		IoDeclarationStatementContext _localctx = new IoDeclarationStatementContext(_ctx, getState());
		enterRule(_localctx, 56, RULE_ioDeclarationStatement);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(387);
			_la = _input.LA(1);
			if ( !(_la==INPUT || _la==OUTPUT) ) {
			_errHandler.recoverInline(this);
			}
			else {
				if ( _input.LA(1)==Token.EOF ) matchedEOF = true;
				_errHandler.reportMatch(this);
				consume();
			}
			setState(390);
			_errHandler.sync(this);
			switch (_input.LA(1)) {
			case BOOL:
			case BIT:
			case INT:
			case UINT:
			case FLOAT:
			case ANGLE:
			case COMPLEX:
			case DURATION:
			case STRETCH:
				{
				setState(388);
				scalarType();
				}
				break;
			case ARRAY:
				{
				setState(389);
				arrayType();
				}
				break;
			default:
				throw new NoViableAltException(this);
			}
			setState(392);
			match(Identifier);
			setState(393);
			match(SEMICOLON);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class OldStyleDeclarationStatementContext extends ParserRuleContext {
		public TerminalNode Identifier() { return getToken(qasm3Parser.Identifier, 0); }
		public TerminalNode SEMICOLON() { return getToken(qasm3Parser.SEMICOLON, 0); }
		public TerminalNode CREG() { return getToken(qasm3Parser.CREG, 0); }
		public TerminalNode QREG() { return getToken(qasm3Parser.QREG, 0); }
		public DesignatorContext designator() {
			return getRuleContext(DesignatorContext.class,0);
		}
		public OldStyleDeclarationStatementContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_oldStyleDeclarationStatement; }
	}

	public final OldStyleDeclarationStatementContext oldStyleDeclarationStatement() throws RecognitionException {
		OldStyleDeclarationStatementContext _localctx = new OldStyleDeclarationStatementContext(_ctx, getState());
		enterRule(_localctx, 58, RULE_oldStyleDeclarationStatement);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(395);
			_la = _input.LA(1);
			if ( !(_la==QREG || _la==CREG) ) {
			_errHandler.recoverInline(this);
			}
			else {
				if ( _input.LA(1)==Token.EOF ) matchedEOF = true;
				_errHandler.reportMatch(this);
				consume();
			}
			setState(396);
			match(Identifier);
			setState(398);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if (_la==LBRACKET) {
				{
				setState(397);
				designator();
				}
			}

			setState(400);
			match(SEMICOLON);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class QuantumDeclarationStatementContext extends ParserRuleContext {
		public QubitTypeContext qubitType() {
			return getRuleContext(QubitTypeContext.class,0);
		}
		public TerminalNode Identifier() { return getToken(qasm3Parser.Identifier, 0); }
		public TerminalNode SEMICOLON() { return getToken(qasm3Parser.SEMICOLON, 0); }
		public QuantumDeclarationStatementContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_quantumDeclarationStatement; }
	}

	public final QuantumDeclarationStatementContext quantumDeclarationStatement() throws RecognitionException {
		QuantumDeclarationStatementContext _localctx = new QuantumDeclarationStatementContext(_ctx, getState());
		enterRule(_localctx, 60, RULE_quantumDeclarationStatement);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(402);
			qubitType();
			setState(403);
			match(Identifier);
			setState(404);
			match(SEMICOLON);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class DefStatementContext extends ParserRuleContext {
		public TerminalNode DEF() { return getToken(qasm3Parser.DEF, 0); }
		public TerminalNode Identifier() { return getToken(qasm3Parser.Identifier, 0); }
		public TerminalNode LPAREN() { return getToken(qasm3Parser.LPAREN, 0); }
		public TerminalNode RPAREN() { return getToken(qasm3Parser.RPAREN, 0); }
		public ScopeContext scope() {
			return getRuleContext(ScopeContext.class,0);
		}
		public ArgumentDefinitionListContext argumentDefinitionList() {
			return getRuleContext(ArgumentDefinitionListContext.class,0);
		}
		public ReturnSignatureContext returnSignature() {
			return getRuleContext(ReturnSignatureContext.class,0);
		}
		public DefStatementContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_defStatement; }
	}

	public final DefStatementContext defStatement() throws RecognitionException {
		DefStatementContext _localctx = new DefStatementContext(_ctx, getState());
		enterRule(_localctx, 62, RULE_defStatement);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(406);
			match(DEF);
			setState(407);
			match(Identifier);
			setState(408);
			match(LPAREN);
			setState(410);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if ((((_la) & ~0x3f) == 0 && ((1L << _la) & 28586765451264L) != 0)) {
				{
				setState(409);
				argumentDefinitionList();
				}
			}

			setState(412);
			match(RPAREN);
			setState(414);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if (_la==ARROW) {
				{
				setState(413);
				returnSignature();
				}
			}

			setState(416);
			scope();
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class ExternStatementContext extends ParserRuleContext {
		public TerminalNode EXTERN() { return getToken(qasm3Parser.EXTERN, 0); }
		public TerminalNode Identifier() { return getToken(qasm3Parser.Identifier, 0); }
		public TerminalNode LPAREN() { return getToken(qasm3Parser.LPAREN, 0); }
		public TerminalNode RPAREN() { return getToken(qasm3Parser.RPAREN, 0); }
		public TerminalNode SEMICOLON() { return getToken(qasm3Parser.SEMICOLON, 0); }
		public ExternArgumentListContext externArgumentList() {
			return getRuleContext(ExternArgumentListContext.class,0);
		}
		public ReturnSignatureContext returnSignature() {
			return getRuleContext(ReturnSignatureContext.class,0);
		}
		public ExternStatementContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_externStatement; }
	}

	public final ExternStatementContext externStatement() throws RecognitionException {
		ExternStatementContext _localctx = new ExternStatementContext(_ctx, getState());
		enterRule(_localctx, 64, RULE_externStatement);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(418);
			match(EXTERN);
			setState(419);
			match(Identifier);
			setState(420);
			match(LPAREN);
			setState(422);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if ((((_la) & ~0x3f) == 0 && ((1L << _la) & 28580323000320L) != 0)) {
				{
				setState(421);
				externArgumentList();
				}
			}

			setState(424);
			match(RPAREN);
			setState(426);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if (_la==ARROW) {
				{
				setState(425);
				returnSignature();
				}
			}

			setState(428);
			match(SEMICOLON);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class GateStatementContext extends ParserRuleContext {
		public IdentifierListContext params;
		public IdentifierListContext qubits;
		public TerminalNode GATE() { return getToken(qasm3Parser.GATE, 0); }
		public TerminalNode Identifier() { return getToken(qasm3Parser.Identifier, 0); }
		public ScopeContext scope() {
			return getRuleContext(ScopeContext.class,0);
		}
		public List<IdentifierListContext> identifierList() {
			return getRuleContexts(IdentifierListContext.class);
		}
		public IdentifierListContext identifierList(int i) {
			return getRuleContext(IdentifierListContext.class,i);
		}
		public TerminalNode LPAREN() { return getToken(qasm3Parser.LPAREN, 0); }
		public TerminalNode RPAREN() { return getToken(qasm3Parser.RPAREN, 0); }
		public GateStatementContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_gateStatement; }
	}

	public final GateStatementContext gateStatement() throws RecognitionException {
		GateStatementContext _localctx = new GateStatementContext(_ctx, getState());
		enterRule(_localctx, 66, RULE_gateStatement);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(430);
			match(GATE);
			setState(431);
			match(Identifier);
			setState(437);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if (_la==LPAREN) {
				{
				setState(432);
				match(LPAREN);
				setState(434);
				_errHandler.sync(this);
				_la = _input.LA(1);
				if (_la==Identifier) {
					{
					setState(433);
					((GateStatementContext)_localctx).params = identifierList();
					}
				}

				setState(436);
				match(RPAREN);
				}
			}

			setState(439);
			((GateStatementContext)_localctx).qubits = identifierList();
			setState(440);
			scope();
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class AssignmentStatementContext extends ParserRuleContext {
		public Token op;
		public IndexedIdentifierContext indexedIdentifier() {
			return getRuleContext(IndexedIdentifierContext.class,0);
		}
		public TerminalNode SEMICOLON() { return getToken(qasm3Parser.SEMICOLON, 0); }
		public TerminalNode EQUALS() { return getToken(qasm3Parser.EQUALS, 0); }
		public TerminalNode CompoundAssignmentOperator() { return getToken(qasm3Parser.CompoundAssignmentOperator, 0); }
		public ExpressionContext expression() {
			return getRuleContext(ExpressionContext.class,0);
		}
		public MeasureExpressionContext measureExpression() {
			return getRuleContext(MeasureExpressionContext.class,0);
		}
		public AssignmentStatementContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_assignmentStatement; }
	}

	public final AssignmentStatementContext assignmentStatement() throws RecognitionException {
		AssignmentStatementContext _localctx = new AssignmentStatementContext(_ctx, getState());
		enterRule(_localctx, 68, RULE_assignmentStatement);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(442);
			indexedIdentifier();
			setState(443);
			((AssignmentStatementContext)_localctx).op = _input.LT(1);
			_la = _input.LA(1);
			if ( !(_la==EQUALS || _la==CompoundAssignmentOperator) ) {
				((AssignmentStatementContext)_localctx).op = (Token)_errHandler.recoverInline(this);
			}
			else {
				if ( _input.LA(1)==Token.EOF ) matchedEOF = true;
				_errHandler.reportMatch(this);
				consume();
			}
			setState(446);
			_errHandler.sync(this);
			switch (_input.LA(1)) {
			case BOOL:
			case BIT:
			case INT:
			case UINT:
			case FLOAT:
			case ANGLE:
			case COMPLEX:
			case ARRAY:
			case DURATION:
			case STRETCH:
			case DURATIONOF:
			case BooleanLiteral:
			case LPAREN:
			case MINUS:
			case TILDE:
			case EXCLAMATION_POINT:
			case ImaginaryLiteral:
			case BinaryIntegerLiteral:
			case OctalIntegerLiteral:
			case DecimalIntegerLiteral:
			case HexIntegerLiteral:
			case Identifier:
			case HardwareQubit:
			case FloatLiteral:
			case TimingLiteral:
			case BitstringLiteral:
				{
				setState(444);
				expression(0);
				}
				break;
			case MEASURE:
				{
				setState(445);
				measureExpression();
				}
				break;
			default:
				throw new NoViableAltException(this);
			}
			setState(448);
			match(SEMICOLON);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class ExpressionStatementContext extends ParserRuleContext {
		public ExpressionContext expression() {
			return getRuleContext(ExpressionContext.class,0);
		}
		public TerminalNode SEMICOLON() { return getToken(qasm3Parser.SEMICOLON, 0); }
		public ExpressionStatementContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_expressionStatement; }
	}

	public final ExpressionStatementContext expressionStatement() throws RecognitionException {
		ExpressionStatementContext _localctx = new ExpressionStatementContext(_ctx, getState());
		enterRule(_localctx, 70, RULE_expressionStatement);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(450);
			expression(0);
			setState(451);
			match(SEMICOLON);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class CalStatementContext extends ParserRuleContext {
		public TerminalNode CAL() { return getToken(qasm3Parser.CAL, 0); }
		public TerminalNode LBRACE() { return getToken(qasm3Parser.LBRACE, 0); }
		public TerminalNode RBRACE() { return getToken(qasm3Parser.RBRACE, 0); }
		public TerminalNode CalibrationBlock() { return getToken(qasm3Parser.CalibrationBlock, 0); }
		public CalStatementContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_calStatement; }
	}

	public final CalStatementContext calStatement() throws RecognitionException {
		CalStatementContext _localctx = new CalStatementContext(_ctx, getState());
		enterRule(_localctx, 72, RULE_calStatement);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(453);
			match(CAL);
			setState(454);
			match(LBRACE);
			setState(456);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if (_la==CalibrationBlock) {
				{
				setState(455);
				match(CalibrationBlock);
				}
			}

			setState(458);
			match(RBRACE);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class DefcalStatementContext extends ParserRuleContext {
		public TerminalNode DEFCAL() { return getToken(qasm3Parser.DEFCAL, 0); }
		public DefcalTargetContext defcalTarget() {
			return getRuleContext(DefcalTargetContext.class,0);
		}
		public DefcalOperandListContext defcalOperandList() {
			return getRuleContext(DefcalOperandListContext.class,0);
		}
		public TerminalNode LBRACE() { return getToken(qasm3Parser.LBRACE, 0); }
		public TerminalNode RBRACE() { return getToken(qasm3Parser.RBRACE, 0); }
		public TerminalNode LPAREN() { return getToken(qasm3Parser.LPAREN, 0); }
		public TerminalNode RPAREN() { return getToken(qasm3Parser.RPAREN, 0); }
		public ReturnSignatureContext returnSignature() {
			return getRuleContext(ReturnSignatureContext.class,0);
		}
		public TerminalNode CalibrationBlock() { return getToken(qasm3Parser.CalibrationBlock, 0); }
		public DefcalArgumentDefinitionListContext defcalArgumentDefinitionList() {
			return getRuleContext(DefcalArgumentDefinitionListContext.class,0);
		}
		public DefcalStatementContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_defcalStatement; }
	}

	public final DefcalStatementContext defcalStatement() throws RecognitionException {
		DefcalStatementContext _localctx = new DefcalStatementContext(_ctx, getState());
		enterRule(_localctx, 74, RULE_defcalStatement);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(460);
			match(DEFCAL);
			setState(461);
			defcalTarget();
			setState(467);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if (_la==LPAREN) {
				{
				setState(462);
				match(LPAREN);
				setState(464);
				_errHandler.sync(this);
				_la = _input.LA(1);
				if ((((_la) & ~0x3f) == 0 && ((1L << _la) & 2380183188854013952L) != 0) || ((((_la - 71)) & ~0x3f) == 0 && ((1L << (_la - 71)) & 268179457L) != 0)) {
					{
					setState(463);
					defcalArgumentDefinitionList();
					}
				}

				setState(466);
				match(RPAREN);
				}
			}

			setState(469);
			defcalOperandList();
			setState(471);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if (_la==ARROW) {
				{
				setState(470);
				returnSignature();
				}
			}

			setState(473);
			match(LBRACE);
			setState(475);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if (_la==CalibrationBlock) {
				{
				setState(474);
				match(CalibrationBlock);
				}
			}

			setState(477);
			match(RBRACE);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class ExpressionContext extends ParserRuleContext {
		public ExpressionContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_expression; }
	 
		public ExpressionContext() { }
		public void copyFrom(ExpressionContext ctx) {
			super.copyFrom(ctx);
		}
	}
	@SuppressWarnings("CheckReturnValue")
	public static class BitwiseXorExpressionContext extends ExpressionContext {
		public Token op;
		public List<ExpressionContext> expression() {
			return getRuleContexts(ExpressionContext.class);
		}
		public ExpressionContext expression(int i) {
			return getRuleContext(ExpressionContext.class,i);
		}
		public TerminalNode CARET() { return getToken(qasm3Parser.CARET, 0); }
		public BitwiseXorExpressionContext(ExpressionContext ctx) { copyFrom(ctx); }
	}
	@SuppressWarnings("CheckReturnValue")
	public static class AdditiveExpressionContext extends ExpressionContext {
		public Token op;
		public List<ExpressionContext> expression() {
			return getRuleContexts(ExpressionContext.class);
		}
		public ExpressionContext expression(int i) {
			return getRuleContext(ExpressionContext.class,i);
		}
		public TerminalNode PLUS() { return getToken(qasm3Parser.PLUS, 0); }
		public TerminalNode MINUS() { return getToken(qasm3Parser.MINUS, 0); }
		public AdditiveExpressionContext(ExpressionContext ctx) { copyFrom(ctx); }
	}
	@SuppressWarnings("CheckReturnValue")
	public static class DurationofExpressionContext extends ExpressionContext {
		public TerminalNode DURATIONOF() { return getToken(qasm3Parser.DURATIONOF, 0); }
		public TerminalNode LPAREN() { return getToken(qasm3Parser.LPAREN, 0); }
		public ScopeContext scope() {
			return getRuleContext(ScopeContext.class,0);
		}
		public TerminalNode RPAREN() { return getToken(qasm3Parser.RPAREN, 0); }
		public DurationofExpressionContext(ExpressionContext ctx) { copyFrom(ctx); }
	}
	@SuppressWarnings("CheckReturnValue")
	public static class ParenthesisExpressionContext extends ExpressionContext {
		public TerminalNode LPAREN() { return getToken(qasm3Parser.LPAREN, 0); }
		public ExpressionContext expression() {
			return getRuleContext(ExpressionContext.class,0);
		}
		public TerminalNode RPAREN() { return getToken(qasm3Parser.RPAREN, 0); }
		public ParenthesisExpressionContext(ExpressionContext ctx) { copyFrom(ctx); }
	}
	@SuppressWarnings("CheckReturnValue")
	public static class ComparisonExpressionContext extends ExpressionContext {
		public Token op;
		public List<ExpressionContext> expression() {
			return getRuleContexts(ExpressionContext.class);
		}
		public ExpressionContext expression(int i) {
			return getRuleContext(ExpressionContext.class,i);
		}
		public TerminalNode ComparisonOperator() { return getToken(qasm3Parser.ComparisonOperator, 0); }
		public ComparisonExpressionContext(ExpressionContext ctx) { copyFrom(ctx); }
	}
	@SuppressWarnings("CheckReturnValue")
	public static class MultiplicativeExpressionContext extends ExpressionContext {
		public Token op;
		public List<ExpressionContext> expression() {
			return getRuleContexts(ExpressionContext.class);
		}
		public ExpressionContext expression(int i) {
			return getRuleContext(ExpressionContext.class,i);
		}
		public TerminalNode ASTERISK() { return getToken(qasm3Parser.ASTERISK, 0); }
		public TerminalNode SLASH() { return getToken(qasm3Parser.SLASH, 0); }
		public TerminalNode PERCENT() { return getToken(qasm3Parser.PERCENT, 0); }
		public MultiplicativeExpressionContext(ExpressionContext ctx) { copyFrom(ctx); }
	}
	@SuppressWarnings("CheckReturnValue")
	public static class LogicalOrExpressionContext extends ExpressionContext {
		public Token op;
		public List<ExpressionContext> expression() {
			return getRuleContexts(ExpressionContext.class);
		}
		public ExpressionContext expression(int i) {
			return getRuleContext(ExpressionContext.class,i);
		}
		public TerminalNode DOUBLE_PIPE() { return getToken(qasm3Parser.DOUBLE_PIPE, 0); }
		public LogicalOrExpressionContext(ExpressionContext ctx) { copyFrom(ctx); }
	}
	@SuppressWarnings("CheckReturnValue")
	public static class CastExpressionContext extends ExpressionContext {
		public TerminalNode LPAREN() { return getToken(qasm3Parser.LPAREN, 0); }
		public ExpressionContext expression() {
			return getRuleContext(ExpressionContext.class,0);
		}
		public TerminalNode RPAREN() { return getToken(qasm3Parser.RPAREN, 0); }
		public ScalarTypeContext scalarType() {
			return getRuleContext(ScalarTypeContext.class,0);
		}
		public ArrayTypeContext arrayType() {
			return getRuleContext(ArrayTypeContext.class,0);
		}
		public CastExpressionContext(ExpressionContext ctx) { copyFrom(ctx); }
	}
	@SuppressWarnings("CheckReturnValue")
	public static class PowerExpressionContext extends ExpressionContext {
		public Token op;
		public List<ExpressionContext> expression() {
			return getRuleContexts(ExpressionContext.class);
		}
		public ExpressionContext expression(int i) {
			return getRuleContext(ExpressionContext.class,i);
		}
		public TerminalNode DOUBLE_ASTERISK() { return getToken(qasm3Parser.DOUBLE_ASTERISK, 0); }
		public PowerExpressionContext(ExpressionContext ctx) { copyFrom(ctx); }
	}
	@SuppressWarnings("CheckReturnValue")
	public static class BitwiseOrExpressionContext extends ExpressionContext {
		public Token op;
		public List<ExpressionContext> expression() {
			return getRuleContexts(ExpressionContext.class);
		}
		public ExpressionContext expression(int i) {
			return getRuleContext(ExpressionContext.class,i);
		}
		public TerminalNode PIPE() { return getToken(qasm3Parser.PIPE, 0); }
		public BitwiseOrExpressionContext(ExpressionContext ctx) { copyFrom(ctx); }
	}
	@SuppressWarnings("CheckReturnValue")
	public static class CallExpressionContext extends ExpressionContext {
		public TerminalNode Identifier() { return getToken(qasm3Parser.Identifier, 0); }
		public TerminalNode LPAREN() { return getToken(qasm3Parser.LPAREN, 0); }
		public TerminalNode RPAREN() { return getToken(qasm3Parser.RPAREN, 0); }
		public ExpressionListContext expressionList() {
			return getRuleContext(ExpressionListContext.class,0);
		}
		public CallExpressionContext(ExpressionContext ctx) { copyFrom(ctx); }
	}
	@SuppressWarnings("CheckReturnValue")
	public static class BitshiftExpressionContext extends ExpressionContext {
		public Token op;
		public List<ExpressionContext> expression() {
			return getRuleContexts(ExpressionContext.class);
		}
		public ExpressionContext expression(int i) {
			return getRuleContext(ExpressionContext.class,i);
		}
		public TerminalNode BitshiftOperator() { return getToken(qasm3Parser.BitshiftOperator, 0); }
		public BitshiftExpressionContext(ExpressionContext ctx) { copyFrom(ctx); }
	}
	@SuppressWarnings("CheckReturnValue")
	public static class BitwiseAndExpressionContext extends ExpressionContext {
		public Token op;
		public List<ExpressionContext> expression() {
			return getRuleContexts(ExpressionContext.class);
		}
		public ExpressionContext expression(int i) {
			return getRuleContext(ExpressionContext.class,i);
		}
		public TerminalNode AMPERSAND() { return getToken(qasm3Parser.AMPERSAND, 0); }
		public BitwiseAndExpressionContext(ExpressionContext ctx) { copyFrom(ctx); }
	}
	@SuppressWarnings("CheckReturnValue")
	public static class EqualityExpressionContext extends ExpressionContext {
		public Token op;
		public List<ExpressionContext> expression() {
			return getRuleContexts(ExpressionContext.class);
		}
		public ExpressionContext expression(int i) {
			return getRuleContext(ExpressionContext.class,i);
		}
		public TerminalNode EqualityOperator() { return getToken(qasm3Parser.EqualityOperator, 0); }
		public EqualityExpressionContext(ExpressionContext ctx) { copyFrom(ctx); }
	}
	@SuppressWarnings("CheckReturnValue")
	public static class LogicalAndExpressionContext extends ExpressionContext {
		public Token op;
		public List<ExpressionContext> expression() {
			return getRuleContexts(ExpressionContext.class);
		}
		public ExpressionContext expression(int i) {
			return getRuleContext(ExpressionContext.class,i);
		}
		public TerminalNode DOUBLE_AMPERSAND() { return getToken(qasm3Parser.DOUBLE_AMPERSAND, 0); }
		public LogicalAndExpressionContext(ExpressionContext ctx) { copyFrom(ctx); }
	}
	@SuppressWarnings("CheckReturnValue")
	public static class IndexExpressionContext extends ExpressionContext {
		public ExpressionContext expression() {
			return getRuleContext(ExpressionContext.class,0);
		}
		public IndexOperatorContext indexOperator() {
			return getRuleContext(IndexOperatorContext.class,0);
		}
		public IndexExpressionContext(ExpressionContext ctx) { copyFrom(ctx); }
	}
	@SuppressWarnings("CheckReturnValue")
	public static class UnaryExpressionContext extends ExpressionContext {
		public Token op;
		public ExpressionContext expression() {
			return getRuleContext(ExpressionContext.class,0);
		}
		public TerminalNode TILDE() { return getToken(qasm3Parser.TILDE, 0); }
		public TerminalNode EXCLAMATION_POINT() { return getToken(qasm3Parser.EXCLAMATION_POINT, 0); }
		public TerminalNode MINUS() { return getToken(qasm3Parser.MINUS, 0); }
		public UnaryExpressionContext(ExpressionContext ctx) { copyFrom(ctx); }
	}
	@SuppressWarnings("CheckReturnValue")
	public static class LiteralExpressionContext extends ExpressionContext {
		public TerminalNode Identifier() { return getToken(qasm3Parser.Identifier, 0); }
		public TerminalNode BinaryIntegerLiteral() { return getToken(qasm3Parser.BinaryIntegerLiteral, 0); }
		public TerminalNode OctalIntegerLiteral() { return getToken(qasm3Parser.OctalIntegerLiteral, 0); }
		public TerminalNode DecimalIntegerLiteral() { return getToken(qasm3Parser.DecimalIntegerLiteral, 0); }
		public TerminalNode HexIntegerLiteral() { return getToken(qasm3Parser.HexIntegerLiteral, 0); }
		public TerminalNode FloatLiteral() { return getToken(qasm3Parser.FloatLiteral, 0); }
		public TerminalNode ImaginaryLiteral() { return getToken(qasm3Parser.ImaginaryLiteral, 0); }
		public TerminalNode BooleanLiteral() { return getToken(qasm3Parser.BooleanLiteral, 0); }
		public TerminalNode BitstringLiteral() { return getToken(qasm3Parser.BitstringLiteral, 0); }
		public TerminalNode TimingLiteral() { return getToken(qasm3Parser.TimingLiteral, 0); }
		public TerminalNode HardwareQubit() { return getToken(qasm3Parser.HardwareQubit, 0); }
		public LiteralExpressionContext(ExpressionContext ctx) { copyFrom(ctx); }
	}

	public final ExpressionContext expression() throws RecognitionException {
		return expression(0);
	}

	private ExpressionContext expression(int _p) throws RecognitionException {
		ParserRuleContext _parentctx = _ctx;
		int _parentState = getState();
		ExpressionContext _localctx = new ExpressionContext(_ctx, _parentState);
		ExpressionContext _prevctx = _localctx;
		int _startState = 76;
		enterRecursionRule(_localctx, 76, RULE_expression, _p);
		int _la;
		try {
			int _alt;
			enterOuterAlt(_localctx, 1);
			{
			setState(506);
			_errHandler.sync(this);
			switch ( getInterpreter().adaptivePredict(_input,46,_ctx) ) {
			case 1:
				{
				_localctx = new ParenthesisExpressionContext(_localctx);
				_ctx = _localctx;
				_prevctx = _localctx;

				setState(480);
				match(LPAREN);
				setState(481);
				expression(0);
				setState(482);
				match(RPAREN);
				}
				break;
			case 2:
				{
				_localctx = new UnaryExpressionContext(_localctx);
				_ctx = _localctx;
				_prevctx = _localctx;
				setState(484);
				((UnaryExpressionContext)_localctx).op = _input.LT(1);
				_la = _input.LA(1);
				if ( !(((((_la - 71)) & ~0x3f) == 0 && ((1L << (_la - 71)) & 6145L) != 0)) ) {
					((UnaryExpressionContext)_localctx).op = (Token)_errHandler.recoverInline(this);
				}
				else {
					if ( _input.LA(1)==Token.EOF ) matchedEOF = true;
					_errHandler.reportMatch(this);
					consume();
				}
				setState(485);
				expression(15);
				}
				break;
			case 3:
				{
				_localctx = new CastExpressionContext(_localctx);
				_ctx = _localctx;
				_prevctx = _localctx;
				setState(488);
				_errHandler.sync(this);
				switch (_input.LA(1)) {
				case BOOL:
				case BIT:
				case INT:
				case UINT:
				case FLOAT:
				case ANGLE:
				case COMPLEX:
				case DURATION:
				case STRETCH:
					{
					setState(486);
					scalarType();
					}
					break;
				case ARRAY:
					{
					setState(487);
					arrayType();
					}
					break;
				default:
					throw new NoViableAltException(this);
				}
				setState(490);
				match(LPAREN);
				setState(491);
				expression(0);
				setState(492);
				match(RPAREN);
				}
				break;
			case 4:
				{
				_localctx = new DurationofExpressionContext(_localctx);
				_ctx = _localctx;
				_prevctx = _localctx;
				setState(494);
				match(DURATIONOF);
				setState(495);
				match(LPAREN);
				setState(496);
				scope();
				setState(497);
				match(RPAREN);
				}
				break;
			case 5:
				{
				_localctx = new CallExpressionContext(_localctx);
				_ctx = _localctx;
				_prevctx = _localctx;
				setState(499);
				match(Identifier);
				setState(500);
				match(LPAREN);
				setState(502);
				_errHandler.sync(this);
				_la = _input.LA(1);
				if ((((_la) & ~0x3f) == 0 && ((1L << _la) & 2380183172211015680L) != 0) || ((((_la - 71)) & ~0x3f) == 0 && ((1L << (_la - 71)) & 268179457L) != 0)) {
					{
					setState(501);
					expressionList();
					}
				}

				setState(504);
				match(RPAREN);
				}
				break;
			case 6:
				{
				_localctx = new LiteralExpressionContext(_localctx);
				_ctx = _localctx;
				_prevctx = _localctx;
				setState(505);
				_la = _input.LA(1);
				if ( !(((((_la - 56)) & ~0x3f) == 0 && ((1L << (_la - 56)) & 8787503087617L) != 0)) ) {
				_errHandler.recoverInline(this);
				}
				else {
					if ( _input.LA(1)==Token.EOF ) matchedEOF = true;
					_errHandler.reportMatch(this);
					consume();
				}
				}
				break;
			}
			_ctx.stop = _input.LT(-1);
			setState(545);
			_errHandler.sync(this);
			_alt = getInterpreter().adaptivePredict(_input,48,_ctx);
			while ( _alt!=2 && _alt!=org.antlr.v4.runtime.atn.ATN.INVALID_ALT_NUMBER ) {
				if ( _alt==1 ) {
					if ( _parseListeners!=null ) triggerExitRuleEvent();
					_prevctx = _localctx;
					{
					setState(543);
					_errHandler.sync(this);
					switch ( getInterpreter().adaptivePredict(_input,47,_ctx) ) {
					case 1:
						{
						_localctx = new PowerExpressionContext(new ExpressionContext(_parentctx, _parentState));
						pushNewRecursionContext(_localctx, _startState, RULE_expression);
						setState(508);
						if (!(precpred(_ctx, 16))) throw new FailedPredicateException(this, "precpred(_ctx, 16)");
						setState(509);
						((PowerExpressionContext)_localctx).op = match(DOUBLE_ASTERISK);
						setState(510);
						expression(16);
						}
						break;
					case 2:
						{
						_localctx = new MultiplicativeExpressionContext(new ExpressionContext(_parentctx, _parentState));
						pushNewRecursionContext(_localctx, _startState, RULE_expression);
						setState(511);
						if (!(precpred(_ctx, 14))) throw new FailedPredicateException(this, "precpred(_ctx, 14)");
						setState(512);
						((MultiplicativeExpressionContext)_localctx).op = _input.LT(1);
						_la = _input.LA(1);
						if ( !(((((_la - 72)) & ~0x3f) == 0 && ((1L << (_la - 72)) & 13L) != 0)) ) {
							((MultiplicativeExpressionContext)_localctx).op = (Token)_errHandler.recoverInline(this);
						}
						else {
							if ( _input.LA(1)==Token.EOF ) matchedEOF = true;
							_errHandler.reportMatch(this);
							consume();
						}
						setState(513);
						expression(15);
						}
						break;
					case 3:
						{
						_localctx = new AdditiveExpressionContext(new ExpressionContext(_parentctx, _parentState));
						pushNewRecursionContext(_localctx, _startState, RULE_expression);
						setState(514);
						if (!(precpred(_ctx, 13))) throw new FailedPredicateException(this, "precpred(_ctx, 13)");
						setState(515);
						((AdditiveExpressionContext)_localctx).op = _input.LT(1);
						_la = _input.LA(1);
						if ( !(_la==PLUS || _la==MINUS) ) {
							((AdditiveExpressionContext)_localctx).op = (Token)_errHandler.recoverInline(this);
						}
						else {
							if ( _input.LA(1)==Token.EOF ) matchedEOF = true;
							_errHandler.reportMatch(this);
							consume();
						}
						setState(516);
						expression(14);
						}
						break;
					case 4:
						{
						_localctx = new BitshiftExpressionContext(new ExpressionContext(_parentctx, _parentState));
						pushNewRecursionContext(_localctx, _startState, RULE_expression);
						setState(517);
						if (!(precpred(_ctx, 12))) throw new FailedPredicateException(this, "precpred(_ctx, 12)");
						setState(518);
						((BitshiftExpressionContext)_localctx).op = match(BitshiftOperator);
						setState(519);
						expression(13);
						}
						break;
					case 5:
						{
						_localctx = new ComparisonExpressionContext(new ExpressionContext(_parentctx, _parentState));
						pushNewRecursionContext(_localctx, _startState, RULE_expression);
						setState(520);
						if (!(precpred(_ctx, 11))) throw new FailedPredicateException(this, "precpred(_ctx, 11)");
						setState(521);
						((ComparisonExpressionContext)_localctx).op = match(ComparisonOperator);
						setState(522);
						expression(12);
						}
						break;
					case 6:
						{
						_localctx = new EqualityExpressionContext(new ExpressionContext(_parentctx, _parentState));
						pushNewRecursionContext(_localctx, _startState, RULE_expression);
						setState(523);
						if (!(precpred(_ctx, 10))) throw new FailedPredicateException(this, "precpred(_ctx, 10)");
						setState(524);
						((EqualityExpressionContext)_localctx).op = match(EqualityOperator);
						setState(525);
						expression(11);
						}
						break;
					case 7:
						{
						_localctx = new BitwiseAndExpressionContext(new ExpressionContext(_parentctx, _parentState));
						pushNewRecursionContext(_localctx, _startState, RULE_expression);
						setState(526);
						if (!(precpred(_ctx, 9))) throw new FailedPredicateException(this, "precpred(_ctx, 9)");
						setState(527);
						((BitwiseAndExpressionContext)_localctx).op = match(AMPERSAND);
						setState(528);
						expression(10);
						}
						break;
					case 8:
						{
						_localctx = new BitwiseXorExpressionContext(new ExpressionContext(_parentctx, _parentState));
						pushNewRecursionContext(_localctx, _startState, RULE_expression);
						setState(529);
						if (!(precpred(_ctx, 8))) throw new FailedPredicateException(this, "precpred(_ctx, 8)");
						setState(530);
						((BitwiseXorExpressionContext)_localctx).op = match(CARET);
						setState(531);
						expression(9);
						}
						break;
					case 9:
						{
						_localctx = new BitwiseOrExpressionContext(new ExpressionContext(_parentctx, _parentState));
						pushNewRecursionContext(_localctx, _startState, RULE_expression);
						setState(532);
						if (!(precpred(_ctx, 7))) throw new FailedPredicateException(this, "precpred(_ctx, 7)");
						setState(533);
						((BitwiseOrExpressionContext)_localctx).op = match(PIPE);
						setState(534);
						expression(8);
						}
						break;
					case 10:
						{
						_localctx = new LogicalAndExpressionContext(new ExpressionContext(_parentctx, _parentState));
						pushNewRecursionContext(_localctx, _startState, RULE_expression);
						setState(535);
						if (!(precpred(_ctx, 6))) throw new FailedPredicateException(this, "precpred(_ctx, 6)");
						setState(536);
						((LogicalAndExpressionContext)_localctx).op = match(DOUBLE_AMPERSAND);
						setState(537);
						expression(7);
						}
						break;
					case 11:
						{
						_localctx = new LogicalOrExpressionContext(new ExpressionContext(_parentctx, _parentState));
						pushNewRecursionContext(_localctx, _startState, RULE_expression);
						setState(538);
						if (!(precpred(_ctx, 5))) throw new FailedPredicateException(this, "precpred(_ctx, 5)");
						setState(539);
						((LogicalOrExpressionContext)_localctx).op = match(DOUBLE_PIPE);
						setState(540);
						expression(6);
						}
						break;
					case 12:
						{
						_localctx = new IndexExpressionContext(new ExpressionContext(_parentctx, _parentState));
						pushNewRecursionContext(_localctx, _startState, RULE_expression);
						setState(541);
						if (!(precpred(_ctx, 17))) throw new FailedPredicateException(this, "precpred(_ctx, 17)");
						setState(542);
						indexOperator();
						}
						break;
					}
					} 
				}
				setState(547);
				_errHandler.sync(this);
				_alt = getInterpreter().adaptivePredict(_input,48,_ctx);
			}
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			unrollRecursionContexts(_parentctx);
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class AliasExpressionContext extends ParserRuleContext {
		public List<ExpressionContext> expression() {
			return getRuleContexts(ExpressionContext.class);
		}
		public ExpressionContext expression(int i) {
			return getRuleContext(ExpressionContext.class,i);
		}
		public List<TerminalNode> DOUBLE_PLUS() { return getTokens(qasm3Parser.DOUBLE_PLUS); }
		public TerminalNode DOUBLE_PLUS(int i) {
			return getToken(qasm3Parser.DOUBLE_PLUS, i);
		}
		public AliasExpressionContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_aliasExpression; }
	}

	public final AliasExpressionContext aliasExpression() throws RecognitionException {
		AliasExpressionContext _localctx = new AliasExpressionContext(_ctx, getState());
		enterRule(_localctx, 78, RULE_aliasExpression);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(548);
			expression(0);
			setState(553);
			_errHandler.sync(this);
			_la = _input.LA(1);
			while (_la==DOUBLE_PLUS) {
				{
				{
				setState(549);
				match(DOUBLE_PLUS);
				setState(550);
				expression(0);
				}
				}
				setState(555);
				_errHandler.sync(this);
				_la = _input.LA(1);
			}
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class DeclarationExpressionContext extends ParserRuleContext {
		public ArrayLiteralContext arrayLiteral() {
			return getRuleContext(ArrayLiteralContext.class,0);
		}
		public ExpressionContext expression() {
			return getRuleContext(ExpressionContext.class,0);
		}
		public MeasureExpressionContext measureExpression() {
			return getRuleContext(MeasureExpressionContext.class,0);
		}
		public DeclarationExpressionContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_declarationExpression; }
	}

	public final DeclarationExpressionContext declarationExpression() throws RecognitionException {
		DeclarationExpressionContext _localctx = new DeclarationExpressionContext(_ctx, getState());
		enterRule(_localctx, 80, RULE_declarationExpression);
		try {
			setState(559);
			_errHandler.sync(this);
			switch (_input.LA(1)) {
			case LBRACE:
				enterOuterAlt(_localctx, 1);
				{
				setState(556);
				arrayLiteral();
				}
				break;
			case BOOL:
			case BIT:
			case INT:
			case UINT:
			case FLOAT:
			case ANGLE:
			case COMPLEX:
			case ARRAY:
			case DURATION:
			case STRETCH:
			case DURATIONOF:
			case BooleanLiteral:
			case LPAREN:
			case MINUS:
			case TILDE:
			case EXCLAMATION_POINT:
			case ImaginaryLiteral:
			case BinaryIntegerLiteral:
			case OctalIntegerLiteral:
			case DecimalIntegerLiteral:
			case HexIntegerLiteral:
			case Identifier:
			case HardwareQubit:
			case FloatLiteral:
			case TimingLiteral:
			case BitstringLiteral:
				enterOuterAlt(_localctx, 2);
				{
				setState(557);
				expression(0);
				}
				break;
			case MEASURE:
				enterOuterAlt(_localctx, 3);
				{
				setState(558);
				measureExpression();
				}
				break;
			default:
				throw new NoViableAltException(this);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class MeasureExpressionContext extends ParserRuleContext {
		public TerminalNode MEASURE() { return getToken(qasm3Parser.MEASURE, 0); }
		public GateOperandContext gateOperand() {
			return getRuleContext(GateOperandContext.class,0);
		}
		public MeasureExpressionContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_measureExpression; }
	}

	public final MeasureExpressionContext measureExpression() throws RecognitionException {
		MeasureExpressionContext _localctx = new MeasureExpressionContext(_ctx, getState());
		enterRule(_localctx, 82, RULE_measureExpression);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(561);
			match(MEASURE);
			setState(562);
			gateOperand();
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class RangeExpressionContext extends ParserRuleContext {
		public List<TerminalNode> COLON() { return getTokens(qasm3Parser.COLON); }
		public TerminalNode COLON(int i) {
			return getToken(qasm3Parser.COLON, i);
		}
		public List<ExpressionContext> expression() {
			return getRuleContexts(ExpressionContext.class);
		}
		public ExpressionContext expression(int i) {
			return getRuleContext(ExpressionContext.class,i);
		}
		public RangeExpressionContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_rangeExpression; }
	}

	public final RangeExpressionContext rangeExpression() throws RecognitionException {
		RangeExpressionContext _localctx = new RangeExpressionContext(_ctx, getState());
		enterRule(_localctx, 84, RULE_rangeExpression);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(565);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if ((((_la) & ~0x3f) == 0 && ((1L << _la) & 2380183172211015680L) != 0) || ((((_la - 71)) & ~0x3f) == 0 && ((1L << (_la - 71)) & 268179457L) != 0)) {
				{
				setState(564);
				expression(0);
				}
			}

			setState(567);
			match(COLON);
			setState(569);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if ((((_la) & ~0x3f) == 0 && ((1L << _la) & 2380183172211015680L) != 0) || ((((_la - 71)) & ~0x3f) == 0 && ((1L << (_la - 71)) & 268179457L) != 0)) {
				{
				setState(568);
				expression(0);
				}
			}

			setState(573);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if (_la==COLON) {
				{
				setState(571);
				match(COLON);
				setState(572);
				expression(0);
				}
			}

			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class SetExpressionContext extends ParserRuleContext {
		public TerminalNode LBRACE() { return getToken(qasm3Parser.LBRACE, 0); }
		public List<ExpressionContext> expression() {
			return getRuleContexts(ExpressionContext.class);
		}
		public ExpressionContext expression(int i) {
			return getRuleContext(ExpressionContext.class,i);
		}
		public TerminalNode RBRACE() { return getToken(qasm3Parser.RBRACE, 0); }
		public List<TerminalNode> COMMA() { return getTokens(qasm3Parser.COMMA); }
		public TerminalNode COMMA(int i) {
			return getToken(qasm3Parser.COMMA, i);
		}
		public SetExpressionContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_setExpression; }
	}

	public final SetExpressionContext setExpression() throws RecognitionException {
		SetExpressionContext _localctx = new SetExpressionContext(_ctx, getState());
		enterRule(_localctx, 86, RULE_setExpression);
		int _la;
		try {
			int _alt;
			enterOuterAlt(_localctx, 1);
			{
			setState(575);
			match(LBRACE);
			setState(576);
			expression(0);
			setState(581);
			_errHandler.sync(this);
			_alt = getInterpreter().adaptivePredict(_input,54,_ctx);
			while ( _alt!=2 && _alt!=org.antlr.v4.runtime.atn.ATN.INVALID_ALT_NUMBER ) {
				if ( _alt==1 ) {
					{
					{
					setState(577);
					match(COMMA);
					setState(578);
					expression(0);
					}
					} 
				}
				setState(583);
				_errHandler.sync(this);
				_alt = getInterpreter().adaptivePredict(_input,54,_ctx);
			}
			setState(585);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if (_la==COMMA) {
				{
				setState(584);
				match(COMMA);
				}
			}

			setState(587);
			match(RBRACE);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class ArrayLiteralContext extends ParserRuleContext {
		public TerminalNode LBRACE() { return getToken(qasm3Parser.LBRACE, 0); }
		public TerminalNode RBRACE() { return getToken(qasm3Parser.RBRACE, 0); }
		public List<ExpressionContext> expression() {
			return getRuleContexts(ExpressionContext.class);
		}
		public ExpressionContext expression(int i) {
			return getRuleContext(ExpressionContext.class,i);
		}
		public List<ArrayLiteralContext> arrayLiteral() {
			return getRuleContexts(ArrayLiteralContext.class);
		}
		public ArrayLiteralContext arrayLiteral(int i) {
			return getRuleContext(ArrayLiteralContext.class,i);
		}
		public List<TerminalNode> COMMA() { return getTokens(qasm3Parser.COMMA); }
		public TerminalNode COMMA(int i) {
			return getToken(qasm3Parser.COMMA, i);
		}
		public ArrayLiteralContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_arrayLiteral; }
	}

	public final ArrayLiteralContext arrayLiteral() throws RecognitionException {
		ArrayLiteralContext _localctx = new ArrayLiteralContext(_ctx, getState());
		enterRule(_localctx, 88, RULE_arrayLiteral);
		int _la;
		try {
			int _alt;
			enterOuterAlt(_localctx, 1);
			{
			setState(589);
			match(LBRACE);
			setState(607);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if ((((_la) & ~0x3f) == 0 && ((1L << _la) & 2956643924514439168L) != 0) || ((((_la - 71)) & ~0x3f) == 0 && ((1L << (_la - 71)) & 268179457L) != 0)) {
				{
				setState(592);
				_errHandler.sync(this);
				switch (_input.LA(1)) {
				case BOOL:
				case BIT:
				case INT:
				case UINT:
				case FLOAT:
				case ANGLE:
				case COMPLEX:
				case ARRAY:
				case DURATION:
				case STRETCH:
				case DURATIONOF:
				case BooleanLiteral:
				case LPAREN:
				case MINUS:
				case TILDE:
				case EXCLAMATION_POINT:
				case ImaginaryLiteral:
				case BinaryIntegerLiteral:
				case OctalIntegerLiteral:
				case DecimalIntegerLiteral:
				case HexIntegerLiteral:
				case Identifier:
				case HardwareQubit:
				case FloatLiteral:
				case TimingLiteral:
				case BitstringLiteral:
					{
					setState(590);
					expression(0);
					}
					break;
				case LBRACE:
					{
					setState(591);
					arrayLiteral();
					}
					break;
				default:
					throw new NoViableAltException(this);
				}
				setState(601);
				_errHandler.sync(this);
				_alt = getInterpreter().adaptivePredict(_input,58,_ctx);
				while ( _alt!=2 && _alt!=org.antlr.v4.runtime.atn.ATN.INVALID_ALT_NUMBER ) {
					if ( _alt==1 ) {
						{
						{
						setState(594);
						match(COMMA);
						setState(597);
						_errHandler.sync(this);
						switch (_input.LA(1)) {
						case BOOL:
						case BIT:
						case INT:
						case UINT:
						case FLOAT:
						case ANGLE:
						case COMPLEX:
						case ARRAY:
						case DURATION:
						case STRETCH:
						case DURATIONOF:
						case BooleanLiteral:
						case LPAREN:
						case MINUS:
						case TILDE:
						case EXCLAMATION_POINT:
						case ImaginaryLiteral:
						case BinaryIntegerLiteral:
						case OctalIntegerLiteral:
						case DecimalIntegerLiteral:
						case HexIntegerLiteral:
						case Identifier:
						case HardwareQubit:
						case FloatLiteral:
						case TimingLiteral:
						case BitstringLiteral:
							{
							setState(595);
							expression(0);
							}
							break;
						case LBRACE:
							{
							setState(596);
							arrayLiteral();
							}
							break;
						default:
							throw new NoViableAltException(this);
						}
						}
						} 
					}
					setState(603);
					_errHandler.sync(this);
					_alt = getInterpreter().adaptivePredict(_input,58,_ctx);
				}
				setState(605);
				_errHandler.sync(this);
				_la = _input.LA(1);
				if (_la==COMMA) {
					{
					setState(604);
					match(COMMA);
					}
				}

				}
			}

			setState(609);
			match(RBRACE);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class IndexOperatorContext extends ParserRuleContext {
		public TerminalNode LBRACKET() { return getToken(qasm3Parser.LBRACKET, 0); }
		public TerminalNode RBRACKET() { return getToken(qasm3Parser.RBRACKET, 0); }
		public SetExpressionContext setExpression() {
			return getRuleContext(SetExpressionContext.class,0);
		}
		public List<ExpressionContext> expression() {
			return getRuleContexts(ExpressionContext.class);
		}
		public ExpressionContext expression(int i) {
			return getRuleContext(ExpressionContext.class,i);
		}
		public List<RangeExpressionContext> rangeExpression() {
			return getRuleContexts(RangeExpressionContext.class);
		}
		public RangeExpressionContext rangeExpression(int i) {
			return getRuleContext(RangeExpressionContext.class,i);
		}
		public List<TerminalNode> COMMA() { return getTokens(qasm3Parser.COMMA); }
		public TerminalNode COMMA(int i) {
			return getToken(qasm3Parser.COMMA, i);
		}
		public IndexOperatorContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_indexOperator; }
	}

	public final IndexOperatorContext indexOperator() throws RecognitionException {
		IndexOperatorContext _localctx = new IndexOperatorContext(_ctx, getState());
		enterRule(_localctx, 90, RULE_indexOperator);
		int _la;
		try {
			int _alt;
			enterOuterAlt(_localctx, 1);
			{
			setState(611);
			match(LBRACKET);
			setState(630);
			_errHandler.sync(this);
			switch (_input.LA(1)) {
			case LBRACE:
				{
				setState(612);
				setExpression();
				}
				break;
			case BOOL:
			case BIT:
			case INT:
			case UINT:
			case FLOAT:
			case ANGLE:
			case COMPLEX:
			case ARRAY:
			case DURATION:
			case STRETCH:
			case DURATIONOF:
			case BooleanLiteral:
			case LPAREN:
			case COLON:
			case MINUS:
			case TILDE:
			case EXCLAMATION_POINT:
			case ImaginaryLiteral:
			case BinaryIntegerLiteral:
			case OctalIntegerLiteral:
			case DecimalIntegerLiteral:
			case HexIntegerLiteral:
			case Identifier:
			case HardwareQubit:
			case FloatLiteral:
			case TimingLiteral:
			case BitstringLiteral:
				{
				setState(615);
				_errHandler.sync(this);
				switch ( getInterpreter().adaptivePredict(_input,61,_ctx) ) {
				case 1:
					{
					setState(613);
					expression(0);
					}
					break;
				case 2:
					{
					setState(614);
					rangeExpression();
					}
					break;
				}
				setState(624);
				_errHandler.sync(this);
				_alt = getInterpreter().adaptivePredict(_input,63,_ctx);
				while ( _alt!=2 && _alt!=org.antlr.v4.runtime.atn.ATN.INVALID_ALT_NUMBER ) {
					if ( _alt==1 ) {
						{
						{
						setState(617);
						match(COMMA);
						setState(620);
						_errHandler.sync(this);
						switch ( getInterpreter().adaptivePredict(_input,62,_ctx) ) {
						case 1:
							{
							setState(618);
							expression(0);
							}
							break;
						case 2:
							{
							setState(619);
							rangeExpression();
							}
							break;
						}
						}
						} 
					}
					setState(626);
					_errHandler.sync(this);
					_alt = getInterpreter().adaptivePredict(_input,63,_ctx);
				}
				setState(628);
				_errHandler.sync(this);
				_la = _input.LA(1);
				if (_la==COMMA) {
					{
					setState(627);
					match(COMMA);
					}
				}

				}
				break;
			default:
				throw new NoViableAltException(this);
			}
			setState(632);
			match(RBRACKET);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class IndexedIdentifierContext extends ParserRuleContext {
		public TerminalNode Identifier() { return getToken(qasm3Parser.Identifier, 0); }
		public List<IndexOperatorContext> indexOperator() {
			return getRuleContexts(IndexOperatorContext.class);
		}
		public IndexOperatorContext indexOperator(int i) {
			return getRuleContext(IndexOperatorContext.class,i);
		}
		public IndexedIdentifierContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_indexedIdentifier; }
	}

	public final IndexedIdentifierContext indexedIdentifier() throws RecognitionException {
		IndexedIdentifierContext _localctx = new IndexedIdentifierContext(_ctx, getState());
		enterRule(_localctx, 92, RULE_indexedIdentifier);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(634);
			match(Identifier);
			setState(638);
			_errHandler.sync(this);
			_la = _input.LA(1);
			while (_la==LBRACKET) {
				{
				{
				setState(635);
				indexOperator();
				}
				}
				setState(640);
				_errHandler.sync(this);
				_la = _input.LA(1);
			}
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class ReturnSignatureContext extends ParserRuleContext {
		public TerminalNode ARROW() { return getToken(qasm3Parser.ARROW, 0); }
		public ScalarTypeContext scalarType() {
			return getRuleContext(ScalarTypeContext.class,0);
		}
		public ReturnSignatureContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_returnSignature; }
	}

	public final ReturnSignatureContext returnSignature() throws RecognitionException {
		ReturnSignatureContext _localctx = new ReturnSignatureContext(_ctx, getState());
		enterRule(_localctx, 94, RULE_returnSignature);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(641);
			match(ARROW);
			setState(642);
			scalarType();
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class GateModifierContext extends ParserRuleContext {
		public TerminalNode AT() { return getToken(qasm3Parser.AT, 0); }
		public TerminalNode INV() { return getToken(qasm3Parser.INV, 0); }
		public TerminalNode POW() { return getToken(qasm3Parser.POW, 0); }
		public TerminalNode LPAREN() { return getToken(qasm3Parser.LPAREN, 0); }
		public ExpressionContext expression() {
			return getRuleContext(ExpressionContext.class,0);
		}
		public TerminalNode RPAREN() { return getToken(qasm3Parser.RPAREN, 0); }
		public TerminalNode CTRL() { return getToken(qasm3Parser.CTRL, 0); }
		public TerminalNode NEGCTRL() { return getToken(qasm3Parser.NEGCTRL, 0); }
		public GateModifierContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_gateModifier; }
	}

	public final GateModifierContext gateModifier() throws RecognitionException {
		GateModifierContext _localctx = new GateModifierContext(_ctx, getState());
		enterRule(_localctx, 96, RULE_gateModifier);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(657);
			_errHandler.sync(this);
			switch (_input.LA(1)) {
			case INV:
				{
				setState(644);
				match(INV);
				}
				break;
			case POW:
				{
				setState(645);
				match(POW);
				setState(646);
				match(LPAREN);
				setState(647);
				expression(0);
				setState(648);
				match(RPAREN);
				}
				break;
			case CTRL:
			case NEGCTRL:
				{
				setState(650);
				_la = _input.LA(1);
				if ( !(_la==CTRL || _la==NEGCTRL) ) {
				_errHandler.recoverInline(this);
				}
				else {
					if ( _input.LA(1)==Token.EOF ) matchedEOF = true;
					_errHandler.reportMatch(this);
					consume();
				}
				setState(655);
				_errHandler.sync(this);
				_la = _input.LA(1);
				if (_la==LPAREN) {
					{
					setState(651);
					match(LPAREN);
					setState(652);
					expression(0);
					setState(653);
					match(RPAREN);
					}
				}

				}
				break;
			default:
				throw new NoViableAltException(this);
			}
			setState(659);
			match(AT);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class ScalarTypeContext extends ParserRuleContext {
		public TerminalNode BIT() { return getToken(qasm3Parser.BIT, 0); }
		public DesignatorContext designator() {
			return getRuleContext(DesignatorContext.class,0);
		}
		public TerminalNode INT() { return getToken(qasm3Parser.INT, 0); }
		public TerminalNode UINT() { return getToken(qasm3Parser.UINT, 0); }
		public TerminalNode FLOAT() { return getToken(qasm3Parser.FLOAT, 0); }
		public TerminalNode ANGLE() { return getToken(qasm3Parser.ANGLE, 0); }
		public TerminalNode BOOL() { return getToken(qasm3Parser.BOOL, 0); }
		public TerminalNode DURATION() { return getToken(qasm3Parser.DURATION, 0); }
		public TerminalNode STRETCH() { return getToken(qasm3Parser.STRETCH, 0); }
		public TerminalNode COMPLEX() { return getToken(qasm3Parser.COMPLEX, 0); }
		public TerminalNode LBRACKET() { return getToken(qasm3Parser.LBRACKET, 0); }
		public ScalarTypeContext scalarType() {
			return getRuleContext(ScalarTypeContext.class,0);
		}
		public TerminalNode RBRACKET() { return getToken(qasm3Parser.RBRACKET, 0); }
		public ScalarTypeContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_scalarType; }
	}

	public final ScalarTypeContext scalarType() throws RecognitionException {
		ScalarTypeContext _localctx = new ScalarTypeContext(_ctx, getState());
		enterRule(_localctx, 98, RULE_scalarType);
		int _la;
		try {
			setState(691);
			_errHandler.sync(this);
			switch (_input.LA(1)) {
			case BIT:
				enterOuterAlt(_localctx, 1);
				{
				setState(661);
				match(BIT);
				setState(663);
				_errHandler.sync(this);
				_la = _input.LA(1);
				if (_la==LBRACKET) {
					{
					setState(662);
					designator();
					}
				}

				}
				break;
			case INT:
				enterOuterAlt(_localctx, 2);
				{
				setState(665);
				match(INT);
				setState(667);
				_errHandler.sync(this);
				_la = _input.LA(1);
				if (_la==LBRACKET) {
					{
					setState(666);
					designator();
					}
				}

				}
				break;
			case UINT:
				enterOuterAlt(_localctx, 3);
				{
				setState(669);
				match(UINT);
				setState(671);
				_errHandler.sync(this);
				_la = _input.LA(1);
				if (_la==LBRACKET) {
					{
					setState(670);
					designator();
					}
				}

				}
				break;
			case FLOAT:
				enterOuterAlt(_localctx, 4);
				{
				setState(673);
				match(FLOAT);
				setState(675);
				_errHandler.sync(this);
				_la = _input.LA(1);
				if (_la==LBRACKET) {
					{
					setState(674);
					designator();
					}
				}

				}
				break;
			case ANGLE:
				enterOuterAlt(_localctx, 5);
				{
				setState(677);
				match(ANGLE);
				setState(679);
				_errHandler.sync(this);
				_la = _input.LA(1);
				if (_la==LBRACKET) {
					{
					setState(678);
					designator();
					}
				}

				}
				break;
			case BOOL:
				enterOuterAlt(_localctx, 6);
				{
				setState(681);
				match(BOOL);
				}
				break;
			case DURATION:
				enterOuterAlt(_localctx, 7);
				{
				setState(682);
				match(DURATION);
				}
				break;
			case STRETCH:
				enterOuterAlt(_localctx, 8);
				{
				setState(683);
				match(STRETCH);
				}
				break;
			case COMPLEX:
				enterOuterAlt(_localctx, 9);
				{
				setState(684);
				match(COMPLEX);
				setState(689);
				_errHandler.sync(this);
				_la = _input.LA(1);
				if (_la==LBRACKET) {
					{
					setState(685);
					match(LBRACKET);
					setState(686);
					scalarType();
					setState(687);
					match(RBRACKET);
					}
				}

				}
				break;
			default:
				throw new NoViableAltException(this);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class QubitTypeContext extends ParserRuleContext {
		public TerminalNode QUBIT() { return getToken(qasm3Parser.QUBIT, 0); }
		public DesignatorContext designator() {
			return getRuleContext(DesignatorContext.class,0);
		}
		public QubitTypeContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_qubitType; }
	}

	public final QubitTypeContext qubitType() throws RecognitionException {
		QubitTypeContext _localctx = new QubitTypeContext(_ctx, getState());
		enterRule(_localctx, 100, RULE_qubitType);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(693);
			match(QUBIT);
			setState(695);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if (_la==LBRACKET) {
				{
				setState(694);
				designator();
				}
			}

			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class ArrayTypeContext extends ParserRuleContext {
		public TerminalNode ARRAY() { return getToken(qasm3Parser.ARRAY, 0); }
		public TerminalNode LBRACKET() { return getToken(qasm3Parser.LBRACKET, 0); }
		public ScalarTypeContext scalarType() {
			return getRuleContext(ScalarTypeContext.class,0);
		}
		public TerminalNode COMMA() { return getToken(qasm3Parser.COMMA, 0); }
		public ExpressionListContext expressionList() {
			return getRuleContext(ExpressionListContext.class,0);
		}
		public TerminalNode RBRACKET() { return getToken(qasm3Parser.RBRACKET, 0); }
		public ArrayTypeContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_arrayType; }
	}

	public final ArrayTypeContext arrayType() throws RecognitionException {
		ArrayTypeContext _localctx = new ArrayTypeContext(_ctx, getState());
		enterRule(_localctx, 102, RULE_arrayType);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(697);
			match(ARRAY);
			setState(698);
			match(LBRACKET);
			setState(699);
			scalarType();
			setState(700);
			match(COMMA);
			setState(701);
			expressionList();
			setState(702);
			match(RBRACKET);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class ArrayReferenceTypeContext extends ParserRuleContext {
		public TerminalNode ARRAY() { return getToken(qasm3Parser.ARRAY, 0); }
		public TerminalNode LBRACKET() { return getToken(qasm3Parser.LBRACKET, 0); }
		public ScalarTypeContext scalarType() {
			return getRuleContext(ScalarTypeContext.class,0);
		}
		public TerminalNode COMMA() { return getToken(qasm3Parser.COMMA, 0); }
		public TerminalNode RBRACKET() { return getToken(qasm3Parser.RBRACKET, 0); }
		public TerminalNode READONLY() { return getToken(qasm3Parser.READONLY, 0); }
		public TerminalNode MUTABLE() { return getToken(qasm3Parser.MUTABLE, 0); }
		public ExpressionListContext expressionList() {
			return getRuleContext(ExpressionListContext.class,0);
		}
		public TerminalNode DIM() { return getToken(qasm3Parser.DIM, 0); }
		public TerminalNode EQUALS() { return getToken(qasm3Parser.EQUALS, 0); }
		public ExpressionContext expression() {
			return getRuleContext(ExpressionContext.class,0);
		}
		public ArrayReferenceTypeContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_arrayReferenceType; }
	}

	public final ArrayReferenceTypeContext arrayReferenceType() throws RecognitionException {
		ArrayReferenceTypeContext _localctx = new ArrayReferenceTypeContext(_ctx, getState());
		enterRule(_localctx, 104, RULE_arrayReferenceType);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(704);
			_la = _input.LA(1);
			if ( !(_la==READONLY || _la==MUTABLE) ) {
			_errHandler.recoverInline(this);
			}
			else {
				if ( _input.LA(1)==Token.EOF ) matchedEOF = true;
				_errHandler.reportMatch(this);
				consume();
			}
			setState(705);
			match(ARRAY);
			setState(706);
			match(LBRACKET);
			setState(707);
			scalarType();
			setState(708);
			match(COMMA);
			setState(713);
			_errHandler.sync(this);
			switch (_input.LA(1)) {
			case BOOL:
			case BIT:
			case INT:
			case UINT:
			case FLOAT:
			case ANGLE:
			case COMPLEX:
			case ARRAY:
			case DURATION:
			case STRETCH:
			case DURATIONOF:
			case BooleanLiteral:
			case LPAREN:
			case MINUS:
			case TILDE:
			case EXCLAMATION_POINT:
			case ImaginaryLiteral:
			case BinaryIntegerLiteral:
			case OctalIntegerLiteral:
			case DecimalIntegerLiteral:
			case HexIntegerLiteral:
			case Identifier:
			case HardwareQubit:
			case FloatLiteral:
			case TimingLiteral:
			case BitstringLiteral:
				{
				setState(709);
				expressionList();
				}
				break;
			case DIM:
				{
				setState(710);
				match(DIM);
				setState(711);
				match(EQUALS);
				setState(712);
				expression(0);
				}
				break;
			default:
				throw new NoViableAltException(this);
			}
			setState(715);
			match(RBRACKET);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class DesignatorContext extends ParserRuleContext {
		public TerminalNode LBRACKET() { return getToken(qasm3Parser.LBRACKET, 0); }
		public ExpressionContext expression() {
			return getRuleContext(ExpressionContext.class,0);
		}
		public TerminalNode RBRACKET() { return getToken(qasm3Parser.RBRACKET, 0); }
		public DesignatorContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_designator; }
	}

	public final DesignatorContext designator() throws RecognitionException {
		DesignatorContext _localctx = new DesignatorContext(_ctx, getState());
		enterRule(_localctx, 106, RULE_designator);
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(717);
			match(LBRACKET);
			setState(718);
			expression(0);
			setState(719);
			match(RBRACKET);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class DefcalTargetContext extends ParserRuleContext {
		public TerminalNode MEASURE() { return getToken(qasm3Parser.MEASURE, 0); }
		public TerminalNode RESET() { return getToken(qasm3Parser.RESET, 0); }
		public TerminalNode DELAY() { return getToken(qasm3Parser.DELAY, 0); }
		public TerminalNode Identifier() { return getToken(qasm3Parser.Identifier, 0); }
		public DefcalTargetContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_defcalTarget; }
	}

	public final DefcalTargetContext defcalTarget() throws RecognitionException {
		DefcalTargetContext _localctx = new DefcalTargetContext(_ctx, getState());
		enterRule(_localctx, 108, RULE_defcalTarget);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(721);
			_la = _input.LA(1);
			if ( !(((((_la - 52)) & ~0x3f) == 0 && ((1L << (_la - 52)) & 4398046511111L) != 0)) ) {
			_errHandler.recoverInline(this);
			}
			else {
				if ( _input.LA(1)==Token.EOF ) matchedEOF = true;
				_errHandler.reportMatch(this);
				consume();
			}
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class DefcalArgumentDefinitionContext extends ParserRuleContext {
		public ExpressionContext expression() {
			return getRuleContext(ExpressionContext.class,0);
		}
		public ArgumentDefinitionContext argumentDefinition() {
			return getRuleContext(ArgumentDefinitionContext.class,0);
		}
		public DefcalArgumentDefinitionContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_defcalArgumentDefinition; }
	}

	public final DefcalArgumentDefinitionContext defcalArgumentDefinition() throws RecognitionException {
		DefcalArgumentDefinitionContext _localctx = new DefcalArgumentDefinitionContext(_ctx, getState());
		enterRule(_localctx, 110, RULE_defcalArgumentDefinition);
		try {
			setState(725);
			_errHandler.sync(this);
			switch ( getInterpreter().adaptivePredict(_input,78,_ctx) ) {
			case 1:
				enterOuterAlt(_localctx, 1);
				{
				setState(723);
				expression(0);
				}
				break;
			case 2:
				enterOuterAlt(_localctx, 2);
				{
				setState(724);
				argumentDefinition();
				}
				break;
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class DefcalOperandContext extends ParserRuleContext {
		public TerminalNode HardwareQubit() { return getToken(qasm3Parser.HardwareQubit, 0); }
		public TerminalNode Identifier() { return getToken(qasm3Parser.Identifier, 0); }
		public DefcalOperandContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_defcalOperand; }
	}

	public final DefcalOperandContext defcalOperand() throws RecognitionException {
		DefcalOperandContext _localctx = new DefcalOperandContext(_ctx, getState());
		enterRule(_localctx, 112, RULE_defcalOperand);
		int _la;
		try {
			enterOuterAlt(_localctx, 1);
			{
			setState(727);
			_la = _input.LA(1);
			if ( !(_la==Identifier || _la==HardwareQubit) ) {
			_errHandler.recoverInline(this);
			}
			else {
				if ( _input.LA(1)==Token.EOF ) matchedEOF = true;
				_errHandler.reportMatch(this);
				consume();
			}
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class GateOperandContext extends ParserRuleContext {
		public IndexedIdentifierContext indexedIdentifier() {
			return getRuleContext(IndexedIdentifierContext.class,0);
		}
		public TerminalNode HardwareQubit() { return getToken(qasm3Parser.HardwareQubit, 0); }
		public GateOperandContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_gateOperand; }
	}

	public final GateOperandContext gateOperand() throws RecognitionException {
		GateOperandContext _localctx = new GateOperandContext(_ctx, getState());
		enterRule(_localctx, 114, RULE_gateOperand);
		try {
			setState(731);
			_errHandler.sync(this);
			switch (_input.LA(1)) {
			case Identifier:
				enterOuterAlt(_localctx, 1);
				{
				setState(729);
				indexedIdentifier();
				}
				break;
			case HardwareQubit:
				enterOuterAlt(_localctx, 2);
				{
				setState(730);
				match(HardwareQubit);
				}
				break;
			default:
				throw new NoViableAltException(this);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class ExternArgumentContext extends ParserRuleContext {
		public ScalarTypeContext scalarType() {
			return getRuleContext(ScalarTypeContext.class,0);
		}
		public ArrayReferenceTypeContext arrayReferenceType() {
			return getRuleContext(ArrayReferenceTypeContext.class,0);
		}
		public TerminalNode CREG() { return getToken(qasm3Parser.CREG, 0); }
		public DesignatorContext designator() {
			return getRuleContext(DesignatorContext.class,0);
		}
		public ExternArgumentContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_externArgument; }
	}

	public final ExternArgumentContext externArgument() throws RecognitionException {
		ExternArgumentContext _localctx = new ExternArgumentContext(_ctx, getState());
		enterRule(_localctx, 116, RULE_externArgument);
		int _la;
		try {
			setState(739);
			_errHandler.sync(this);
			switch (_input.LA(1)) {
			case BOOL:
			case BIT:
			case INT:
			case UINT:
			case FLOAT:
			case ANGLE:
			case COMPLEX:
			case DURATION:
			case STRETCH:
				enterOuterAlt(_localctx, 1);
				{
				setState(733);
				scalarType();
				}
				break;
			case READONLY:
			case MUTABLE:
				enterOuterAlt(_localctx, 2);
				{
				setState(734);
				arrayReferenceType();
				}
				break;
			case CREG:
				enterOuterAlt(_localctx, 3);
				{
				setState(735);
				match(CREG);
				setState(737);
				_errHandler.sync(this);
				_la = _input.LA(1);
				if (_la==LBRACKET) {
					{
					setState(736);
					designator();
					}
				}

				}
				break;
			default:
				throw new NoViableAltException(this);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class ArgumentDefinitionContext extends ParserRuleContext {
		public ScalarTypeContext scalarType() {
			return getRuleContext(ScalarTypeContext.class,0);
		}
		public TerminalNode Identifier() { return getToken(qasm3Parser.Identifier, 0); }
		public QubitTypeContext qubitType() {
			return getRuleContext(QubitTypeContext.class,0);
		}
		public TerminalNode CREG() { return getToken(qasm3Parser.CREG, 0); }
		public TerminalNode QREG() { return getToken(qasm3Parser.QREG, 0); }
		public DesignatorContext designator() {
			return getRuleContext(DesignatorContext.class,0);
		}
		public ArrayReferenceTypeContext arrayReferenceType() {
			return getRuleContext(ArrayReferenceTypeContext.class,0);
		}
		public ArgumentDefinitionContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_argumentDefinition; }
	}

	public final ArgumentDefinitionContext argumentDefinition() throws RecognitionException {
		ArgumentDefinitionContext _localctx = new ArgumentDefinitionContext(_ctx, getState());
		enterRule(_localctx, 118, RULE_argumentDefinition);
		int _la;
		try {
			setState(755);
			_errHandler.sync(this);
			switch (_input.LA(1)) {
			case BOOL:
			case BIT:
			case INT:
			case UINT:
			case FLOAT:
			case ANGLE:
			case COMPLEX:
			case DURATION:
			case STRETCH:
				enterOuterAlt(_localctx, 1);
				{
				setState(741);
				scalarType();
				setState(742);
				match(Identifier);
				}
				break;
			case QUBIT:
				enterOuterAlt(_localctx, 2);
				{
				setState(744);
				qubitType();
				setState(745);
				match(Identifier);
				}
				break;
			case QREG:
			case CREG:
				enterOuterAlt(_localctx, 3);
				{
				setState(747);
				_la = _input.LA(1);
				if ( !(_la==QREG || _la==CREG) ) {
				_errHandler.recoverInline(this);
				}
				else {
					if ( _input.LA(1)==Token.EOF ) matchedEOF = true;
					_errHandler.reportMatch(this);
					consume();
				}
				setState(748);
				match(Identifier);
				setState(750);
				_errHandler.sync(this);
				_la = _input.LA(1);
				if (_la==LBRACKET) {
					{
					setState(749);
					designator();
					}
				}

				}
				break;
			case READONLY:
			case MUTABLE:
				enterOuterAlt(_localctx, 4);
				{
				setState(752);
				arrayReferenceType();
				setState(753);
				match(Identifier);
				}
				break;
			default:
				throw new NoViableAltException(this);
			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class ArgumentDefinitionListContext extends ParserRuleContext {
		public List<ArgumentDefinitionContext> argumentDefinition() {
			return getRuleContexts(ArgumentDefinitionContext.class);
		}
		public ArgumentDefinitionContext argumentDefinition(int i) {
			return getRuleContext(ArgumentDefinitionContext.class,i);
		}
		public List<TerminalNode> COMMA() { return getTokens(qasm3Parser.COMMA); }
		public TerminalNode COMMA(int i) {
			return getToken(qasm3Parser.COMMA, i);
		}
		public ArgumentDefinitionListContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_argumentDefinitionList; }
	}

	public final ArgumentDefinitionListContext argumentDefinitionList() throws RecognitionException {
		ArgumentDefinitionListContext _localctx = new ArgumentDefinitionListContext(_ctx, getState());
		enterRule(_localctx, 120, RULE_argumentDefinitionList);
		int _la;
		try {
			int _alt;
			enterOuterAlt(_localctx, 1);
			{
			setState(757);
			argumentDefinition();
			setState(762);
			_errHandler.sync(this);
			_alt = getInterpreter().adaptivePredict(_input,84,_ctx);
			while ( _alt!=2 && _alt!=org.antlr.v4.runtime.atn.ATN.INVALID_ALT_NUMBER ) {
				if ( _alt==1 ) {
					{
					{
					setState(758);
					match(COMMA);
					setState(759);
					argumentDefinition();
					}
					} 
				}
				setState(764);
				_errHandler.sync(this);
				_alt = getInterpreter().adaptivePredict(_input,84,_ctx);
			}
			setState(766);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if (_la==COMMA) {
				{
				setState(765);
				match(COMMA);
				}
			}

			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class DefcalArgumentDefinitionListContext extends ParserRuleContext {
		public List<DefcalArgumentDefinitionContext> defcalArgumentDefinition() {
			return getRuleContexts(DefcalArgumentDefinitionContext.class);
		}
		public DefcalArgumentDefinitionContext defcalArgumentDefinition(int i) {
			return getRuleContext(DefcalArgumentDefinitionContext.class,i);
		}
		public List<TerminalNode> COMMA() { return getTokens(qasm3Parser.COMMA); }
		public TerminalNode COMMA(int i) {
			return getToken(qasm3Parser.COMMA, i);
		}
		public DefcalArgumentDefinitionListContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_defcalArgumentDefinitionList; }
	}

	public final DefcalArgumentDefinitionListContext defcalArgumentDefinitionList() throws RecognitionException {
		DefcalArgumentDefinitionListContext _localctx = new DefcalArgumentDefinitionListContext(_ctx, getState());
		enterRule(_localctx, 122, RULE_defcalArgumentDefinitionList);
		int _la;
		try {
			int _alt;
			enterOuterAlt(_localctx, 1);
			{
			setState(768);
			defcalArgumentDefinition();
			setState(773);
			_errHandler.sync(this);
			_alt = getInterpreter().adaptivePredict(_input,86,_ctx);
			while ( _alt!=2 && _alt!=org.antlr.v4.runtime.atn.ATN.INVALID_ALT_NUMBER ) {
				if ( _alt==1 ) {
					{
					{
					setState(769);
					match(COMMA);
					setState(770);
					defcalArgumentDefinition();
					}
					} 
				}
				setState(775);
				_errHandler.sync(this);
				_alt = getInterpreter().adaptivePredict(_input,86,_ctx);
			}
			setState(777);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if (_la==COMMA) {
				{
				setState(776);
				match(COMMA);
				}
			}

			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class DefcalOperandListContext extends ParserRuleContext {
		public List<DefcalOperandContext> defcalOperand() {
			return getRuleContexts(DefcalOperandContext.class);
		}
		public DefcalOperandContext defcalOperand(int i) {
			return getRuleContext(DefcalOperandContext.class,i);
		}
		public List<TerminalNode> COMMA() { return getTokens(qasm3Parser.COMMA); }
		public TerminalNode COMMA(int i) {
			return getToken(qasm3Parser.COMMA, i);
		}
		public DefcalOperandListContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_defcalOperandList; }
	}

	public final DefcalOperandListContext defcalOperandList() throws RecognitionException {
		DefcalOperandListContext _localctx = new DefcalOperandListContext(_ctx, getState());
		enterRule(_localctx, 124, RULE_defcalOperandList);
		int _la;
		try {
			int _alt;
			enterOuterAlt(_localctx, 1);
			{
			setState(779);
			defcalOperand();
			setState(784);
			_errHandler.sync(this);
			_alt = getInterpreter().adaptivePredict(_input,88,_ctx);
			while ( _alt!=2 && _alt!=org.antlr.v4.runtime.atn.ATN.INVALID_ALT_NUMBER ) {
				if ( _alt==1 ) {
					{
					{
					setState(780);
					match(COMMA);
					setState(781);
					defcalOperand();
					}
					} 
				}
				setState(786);
				_errHandler.sync(this);
				_alt = getInterpreter().adaptivePredict(_input,88,_ctx);
			}
			setState(788);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if (_la==COMMA) {
				{
				setState(787);
				match(COMMA);
				}
			}

			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class ExpressionListContext extends ParserRuleContext {
		public List<ExpressionContext> expression() {
			return getRuleContexts(ExpressionContext.class);
		}
		public ExpressionContext expression(int i) {
			return getRuleContext(ExpressionContext.class,i);
		}
		public List<TerminalNode> COMMA() { return getTokens(qasm3Parser.COMMA); }
		public TerminalNode COMMA(int i) {
			return getToken(qasm3Parser.COMMA, i);
		}
		public ExpressionListContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_expressionList; }
	}

	public final ExpressionListContext expressionList() throws RecognitionException {
		ExpressionListContext _localctx = new ExpressionListContext(_ctx, getState());
		enterRule(_localctx, 126, RULE_expressionList);
		int _la;
		try {
			int _alt;
			enterOuterAlt(_localctx, 1);
			{
			setState(790);
			expression(0);
			setState(795);
			_errHandler.sync(this);
			_alt = getInterpreter().adaptivePredict(_input,90,_ctx);
			while ( _alt!=2 && _alt!=org.antlr.v4.runtime.atn.ATN.INVALID_ALT_NUMBER ) {
				if ( _alt==1 ) {
					{
					{
					setState(791);
					match(COMMA);
					setState(792);
					expression(0);
					}
					} 
				}
				setState(797);
				_errHandler.sync(this);
				_alt = getInterpreter().adaptivePredict(_input,90,_ctx);
			}
			setState(799);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if (_la==COMMA) {
				{
				setState(798);
				match(COMMA);
				}
			}

			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class IdentifierListContext extends ParserRuleContext {
		public List<TerminalNode> Identifier() { return getTokens(qasm3Parser.Identifier); }
		public TerminalNode Identifier(int i) {
			return getToken(qasm3Parser.Identifier, i);
		}
		public List<TerminalNode> COMMA() { return getTokens(qasm3Parser.COMMA); }
		public TerminalNode COMMA(int i) {
			return getToken(qasm3Parser.COMMA, i);
		}
		public IdentifierListContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_identifierList; }
	}

	public final IdentifierListContext identifierList() throws RecognitionException {
		IdentifierListContext _localctx = new IdentifierListContext(_ctx, getState());
		enterRule(_localctx, 128, RULE_identifierList);
		int _la;
		try {
			int _alt;
			enterOuterAlt(_localctx, 1);
			{
			setState(801);
			match(Identifier);
			setState(806);
			_errHandler.sync(this);
			_alt = getInterpreter().adaptivePredict(_input,92,_ctx);
			while ( _alt!=2 && _alt!=org.antlr.v4.runtime.atn.ATN.INVALID_ALT_NUMBER ) {
				if ( _alt==1 ) {
					{
					{
					setState(802);
					match(COMMA);
					setState(803);
					match(Identifier);
					}
					} 
				}
				setState(808);
				_errHandler.sync(this);
				_alt = getInterpreter().adaptivePredict(_input,92,_ctx);
			}
			setState(810);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if (_la==COMMA) {
				{
				setState(809);
				match(COMMA);
				}
			}

			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class GateOperandListContext extends ParserRuleContext {
		public List<GateOperandContext> gateOperand() {
			return getRuleContexts(GateOperandContext.class);
		}
		public GateOperandContext gateOperand(int i) {
			return getRuleContext(GateOperandContext.class,i);
		}
		public List<TerminalNode> COMMA() { return getTokens(qasm3Parser.COMMA); }
		public TerminalNode COMMA(int i) {
			return getToken(qasm3Parser.COMMA, i);
		}
		public GateOperandListContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_gateOperandList; }
	}

	public final GateOperandListContext gateOperandList() throws RecognitionException {
		GateOperandListContext _localctx = new GateOperandListContext(_ctx, getState());
		enterRule(_localctx, 130, RULE_gateOperandList);
		int _la;
		try {
			int _alt;
			enterOuterAlt(_localctx, 1);
			{
			setState(812);
			gateOperand();
			setState(817);
			_errHandler.sync(this);
			_alt = getInterpreter().adaptivePredict(_input,94,_ctx);
			while ( _alt!=2 && _alt!=org.antlr.v4.runtime.atn.ATN.INVALID_ALT_NUMBER ) {
				if ( _alt==1 ) {
					{
					{
					setState(813);
					match(COMMA);
					setState(814);
					gateOperand();
					}
					} 
				}
				setState(819);
				_errHandler.sync(this);
				_alt = getInterpreter().adaptivePredict(_input,94,_ctx);
			}
			setState(821);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if (_la==COMMA) {
				{
				setState(820);
				match(COMMA);
				}
			}

			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	@SuppressWarnings("CheckReturnValue")
	public static class ExternArgumentListContext extends ParserRuleContext {
		public List<ExternArgumentContext> externArgument() {
			return getRuleContexts(ExternArgumentContext.class);
		}
		public ExternArgumentContext externArgument(int i) {
			return getRuleContext(ExternArgumentContext.class,i);
		}
		public List<TerminalNode> COMMA() { return getTokens(qasm3Parser.COMMA); }
		public TerminalNode COMMA(int i) {
			return getToken(qasm3Parser.COMMA, i);
		}
		public ExternArgumentListContext(ParserRuleContext parent, int invokingState) {
			super(parent, invokingState);
		}
		@Override public int getRuleIndex() { return RULE_externArgumentList; }
	}

	public final ExternArgumentListContext externArgumentList() throws RecognitionException {
		ExternArgumentListContext _localctx = new ExternArgumentListContext(_ctx, getState());
		enterRule(_localctx, 132, RULE_externArgumentList);
		int _la;
		try {
			int _alt;
			enterOuterAlt(_localctx, 1);
			{
			setState(823);
			externArgument();
			setState(828);
			_errHandler.sync(this);
			_alt = getInterpreter().adaptivePredict(_input,96,_ctx);
			while ( _alt!=2 && _alt!=org.antlr.v4.runtime.atn.ATN.INVALID_ALT_NUMBER ) {
				if ( _alt==1 ) {
					{
					{
					setState(824);
					match(COMMA);
					setState(825);
					externArgument();
					}
					} 
				}
				setState(830);
				_errHandler.sync(this);
				_alt = getInterpreter().adaptivePredict(_input,96,_ctx);
			}
			setState(832);
			_errHandler.sync(this);
			_la = _input.LA(1);
			if (_la==COMMA) {
				{
				setState(831);
				match(COMMA);
				}
			}

			}
		}
		catch (RecognitionException re) {
			_localctx.exception = re;
			_errHandler.reportError(this, re);
			_errHandler.recover(this, re);
		}
		finally {
			exitRule();
		}
		return _localctx;
	}

	public boolean sempred(RuleContext _localctx, int ruleIndex, int predIndex) {
		switch (ruleIndex) {
		case 38:
			return expression_sempred((ExpressionContext)_localctx, predIndex);
		}
		return true;
	}
	private boolean expression_sempred(ExpressionContext _localctx, int predIndex) {
		switch (predIndex) {
		case 0:
			return precpred(_ctx, 16);
		case 1:
			return precpred(_ctx, 14);
		case 2:
			return precpred(_ctx, 13);
		case 3:
			return precpred(_ctx, 12);
		case 4:
			return precpred(_ctx, 11);
		case 5:
			return precpred(_ctx, 10);
		case 6:
			return precpred(_ctx, 9);
		case 7:
			return precpred(_ctx, 8);
		case 8:
			return precpred(_ctx, 7);
		case 9:
			return precpred(_ctx, 6);
		case 10:
			return precpred(_ctx, 5);
		case 11:
			return precpred(_ctx, 17);
		}
		return true;
	}

	public static final String _serializedATN =
		"\u0004\u0001r\u0343\u0002\u0000\u0007\u0000\u0002\u0001\u0007\u0001\u0002"+
		"\u0002\u0007\u0002\u0002\u0003\u0007\u0003\u0002\u0004\u0007\u0004\u0002"+
		"\u0005\u0007\u0005\u0002\u0006\u0007\u0006\u0002\u0007\u0007\u0007\u0002"+
		"\b\u0007\b\u0002\t\u0007\t\u0002\n\u0007\n\u0002\u000b\u0007\u000b\u0002"+
		"\f\u0007\f\u0002\r\u0007\r\u0002\u000e\u0007\u000e\u0002\u000f\u0007\u000f"+
		"\u0002\u0010\u0007\u0010\u0002\u0011\u0007\u0011\u0002\u0012\u0007\u0012"+
		"\u0002\u0013\u0007\u0013\u0002\u0014\u0007\u0014\u0002\u0015\u0007\u0015"+
		"\u0002\u0016\u0007\u0016\u0002\u0017\u0007\u0017\u0002\u0018\u0007\u0018"+
		"\u0002\u0019\u0007\u0019\u0002\u001a\u0007\u001a\u0002\u001b\u0007\u001b"+
		"\u0002\u001c\u0007\u001c\u0002\u001d\u0007\u001d\u0002\u001e\u0007\u001e"+
		"\u0002\u001f\u0007\u001f\u0002 \u0007 \u0002!\u0007!\u0002\"\u0007\"\u0002"+
		"#\u0007#\u0002$\u0007$\u0002%\u0007%\u0002&\u0007&\u0002\'\u0007\'\u0002"+
		"(\u0007(\u0002)\u0007)\u0002*\u0007*\u0002+\u0007+\u0002,\u0007,\u0002"+
		"-\u0007-\u0002.\u0007.\u0002/\u0007/\u00020\u00070\u00021\u00071\u0002"+
		"2\u00072\u00023\u00073\u00024\u00074\u00025\u00075\u00026\u00076\u0002"+
		"7\u00077\u00028\u00078\u00029\u00079\u0002:\u0007:\u0002;\u0007;\u0002"+
		"<\u0007<\u0002=\u0007=\u0002>\u0007>\u0002?\u0007?\u0002@\u0007@\u0002"+
		"A\u0007A\u0002B\u0007B\u0001\u0000\u0003\u0000\u0088\b\u0000\u0001\u0000"+
		"\u0005\u0000\u008b\b\u0000\n\u0000\f\u0000\u008e\t\u0000\u0001\u0000\u0001"+
		"\u0000\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0002\u0001"+
		"\u0002\u0005\u0002\u0098\b\u0002\n\u0002\f\u0002\u009b\t\u0002\u0001\u0002"+
		"\u0001\u0002\u0001\u0002\u0001\u0002\u0001\u0002\u0001\u0002\u0001\u0002"+
		"\u0001\u0002\u0001\u0002\u0001\u0002\u0001\u0002\u0001\u0002\u0001\u0002"+
		"\u0001\u0002\u0001\u0002\u0001\u0002\u0001\u0002\u0001\u0002\u0001\u0002"+
		"\u0001\u0002\u0001\u0002\u0001\u0002\u0001\u0002\u0001\u0002\u0001\u0002"+
		"\u0001\u0002\u0001\u0002\u0001\u0002\u0001\u0002\u0001\u0002\u0003\u0002"+
		"\u00bb\b\u0002\u0003\u0002\u00bd\b\u0002\u0001\u0003\u0001\u0003\u0003"+
		"\u0003\u00c1\b\u0003\u0001\u0004\u0001\u0004\u0005\u0004\u00c5\b\u0004"+
		"\n\u0004\f\u0004\u00c8\t\u0004\u0001\u0004\u0001\u0004\u0001\u0005\u0001"+
		"\u0005\u0001\u0005\u0001\u0006\u0001\u0006\u0003\u0006\u00d1\b\u0006\u0001"+
		"\u0007\u0001\u0007\u0001\u0007\u0001\u0007\u0001\b\u0001\b\u0001\b\u0001"+
		"\b\u0001\t\u0001\t\u0001\t\u0001\n\u0001\n\u0001\n\u0001\u000b\u0001\u000b"+
		"\u0001\u000b\u0001\f\u0001\f\u0001\f\u0001\f\u0001\f\u0001\f\u0001\f\u0001"+
		"\f\u0001\f\u0001\f\u0003\f\u00ee\b\f\u0001\f\u0001\f\u0001\r\u0001\r\u0001"+
		"\r\u0001\r\u0001\r\u0001\r\u0001\r\u0003\r\u00f9\b\r\u0001\u000e\u0001"+
		"\u000e\u0001\u000e\u0003\u000e\u00fe\b\u000e\u0001\u000e\u0001\u000e\u0001"+
		"\u000f\u0001\u000f\u0001\u000f\u0001\u000f\u0001\u000f\u0001\u000f\u0001"+
		"\u0010\u0001\u0010\u0001\u0010\u0001\u0010\u0001\u0010\u0001\u0010\u0005"+
		"\u0010\u010e\b\u0010\n\u0010\f\u0010\u0111\t\u0010\u0001\u0010\u0001\u0010"+
		"\u0001\u0011\u0001\u0011\u0001\u0011\u0001\u0011\u0001\u0011\u0001\u0011"+
		"\u0003\u0011\u011b\b\u0011\u0001\u0012\u0001\u0012\u0003\u0012\u011f\b"+
		"\u0012\u0001\u0012\u0001\u0012\u0001\u0013\u0001\u0013\u0003\u0013\u0125"+
		"\b\u0013\u0001\u0013\u0001\u0013\u0001\u0014\u0001\u0014\u0001\u0014\u0003"+
		"\u0014\u012c\b\u0014\u0001\u0014\u0001\u0014\u0001\u0015\u0001\u0015\u0003"+
		"\u0015\u0132\b\u0015\u0001\u0015\u0001\u0015\u0001\u0016\u0005\u0016\u0137"+
		"\b\u0016\n\u0016\f\u0016\u013a\t\u0016\u0001\u0016\u0001\u0016\u0001\u0016"+
		"\u0003\u0016\u013f\b\u0016\u0001\u0016\u0003\u0016\u0142\b\u0016\u0001"+
		"\u0016\u0003\u0016\u0145\b\u0016\u0001\u0016\u0001\u0016\u0001\u0016\u0001"+
		"\u0016\u0005\u0016\u014b\b\u0016\n\u0016\f\u0016\u014e\t\u0016\u0001\u0016"+
		"\u0001\u0016\u0001\u0016\u0003\u0016\u0153\b\u0016\u0001\u0016\u0003\u0016"+
		"\u0156\b\u0016\u0001\u0016\u0003\u0016\u0159\b\u0016\u0001\u0016\u0003"+
		"\u0016\u015c\b\u0016\u0001\u0016\u0003\u0016\u015f\b\u0016\u0001\u0017"+
		"\u0001\u0017\u0001\u0017\u0003\u0017\u0164\b\u0017\u0001\u0017\u0001\u0017"+
		"\u0001\u0018\u0001\u0018\u0001\u0018\u0001\u0018\u0001\u0019\u0001\u0019"+
		"\u0001\u0019\u0001\u0019\u0001\u0019\u0001\u0019\u0001\u001a\u0001\u001a"+
		"\u0003\u001a\u0174\b\u001a\u0001\u001a\u0001\u001a\u0001\u001a\u0003\u001a"+
		"\u0179\b\u001a\u0001\u001a\u0001\u001a\u0001\u001b\u0001\u001b\u0001\u001b"+
		"\u0001\u001b\u0001\u001b\u0001\u001b\u0001\u001b\u0001\u001c\u0001\u001c"+
		"\u0001\u001c\u0003\u001c\u0187\b\u001c\u0001\u001c\u0001\u001c\u0001\u001c"+
		"\u0001\u001d\u0001\u001d\u0001\u001d\u0003\u001d\u018f\b\u001d\u0001\u001d"+
		"\u0001\u001d\u0001\u001e\u0001\u001e\u0001\u001e\u0001\u001e\u0001\u001f"+
		"\u0001\u001f\u0001\u001f\u0001\u001f\u0003\u001f\u019b\b\u001f\u0001\u001f"+
		"\u0001\u001f\u0003\u001f\u019f\b\u001f\u0001\u001f\u0001\u001f\u0001 "+
		"\u0001 \u0001 \u0001 \u0003 \u01a7\b \u0001 \u0001 \u0003 \u01ab\b \u0001"+
		" \u0001 \u0001!\u0001!\u0001!\u0001!\u0003!\u01b3\b!\u0001!\u0003!\u01b6"+
		"\b!\u0001!\u0001!\u0001!\u0001\"\u0001\"\u0001\"\u0001\"\u0003\"\u01bf"+
		"\b\"\u0001\"\u0001\"\u0001#\u0001#\u0001#\u0001$\u0001$\u0001$\u0003$"+
		"\u01c9\b$\u0001$\u0001$\u0001%\u0001%\u0001%\u0001%\u0003%\u01d1\b%\u0001"+
		"%\u0003%\u01d4\b%\u0001%\u0001%\u0003%\u01d8\b%\u0001%\u0001%\u0003%\u01dc"+
		"\b%\u0001%\u0001%\u0001&\u0001&\u0001&\u0001&\u0001&\u0001&\u0001&\u0001"+
		"&\u0001&\u0003&\u01e9\b&\u0001&\u0001&\u0001&\u0001&\u0001&\u0001&\u0001"+
		"&\u0001&\u0001&\u0001&\u0001&\u0001&\u0003&\u01f7\b&\u0001&\u0001&\u0003"+
		"&\u01fb\b&\u0001&\u0001&\u0001&\u0001&\u0001&\u0001&\u0001&\u0001&\u0001"+
		"&\u0001&\u0001&\u0001&\u0001&\u0001&\u0001&\u0001&\u0001&\u0001&\u0001"+
		"&\u0001&\u0001&\u0001&\u0001&\u0001&\u0001&\u0001&\u0001&\u0001&\u0001"+
		"&\u0001&\u0001&\u0001&\u0001&\u0001&\u0001&\u0005&\u0220\b&\n&\f&\u0223"+
		"\t&\u0001\'\u0001\'\u0001\'\u0005\'\u0228\b\'\n\'\f\'\u022b\t\'\u0001"+
		"(\u0001(\u0001(\u0003(\u0230\b(\u0001)\u0001)\u0001)\u0001*\u0003*\u0236"+
		"\b*\u0001*\u0001*\u0003*\u023a\b*\u0001*\u0001*\u0003*\u023e\b*\u0001"+
		"+\u0001+\u0001+\u0001+\u0005+\u0244\b+\n+\f+\u0247\t+\u0001+\u0003+\u024a"+
		"\b+\u0001+\u0001+\u0001,\u0001,\u0001,\u0003,\u0251\b,\u0001,\u0001,\u0001"+
		",\u0003,\u0256\b,\u0005,\u0258\b,\n,\f,\u025b\t,\u0001,\u0003,\u025e\b"+
		",\u0003,\u0260\b,\u0001,\u0001,\u0001-\u0001-\u0001-\u0001-\u0003-\u0268"+
		"\b-\u0001-\u0001-\u0001-\u0003-\u026d\b-\u0005-\u026f\b-\n-\f-\u0272\t"+
		"-\u0001-\u0003-\u0275\b-\u0003-\u0277\b-\u0001-\u0001-\u0001.\u0001.\u0005"+
		".\u027d\b.\n.\f.\u0280\t.\u0001/\u0001/\u0001/\u00010\u00010\u00010\u0001"+
		"0\u00010\u00010\u00010\u00010\u00010\u00010\u00010\u00030\u0290\b0\u0003"+
		"0\u0292\b0\u00010\u00010\u00011\u00011\u00031\u0298\b1\u00011\u00011\u0003"+
		"1\u029c\b1\u00011\u00011\u00031\u02a0\b1\u00011\u00011\u00031\u02a4\b"+
		"1\u00011\u00011\u00031\u02a8\b1\u00011\u00011\u00011\u00011\u00011\u0001"+
		"1\u00011\u00011\u00031\u02b2\b1\u00031\u02b4\b1\u00012\u00012\u00032\u02b8"+
		"\b2\u00013\u00013\u00013\u00013\u00013\u00013\u00013\u00014\u00014\u0001"+
		"4\u00014\u00014\u00014\u00014\u00014\u00014\u00034\u02ca\b4\u00014\u0001"+
		"4\u00015\u00015\u00015\u00015\u00016\u00016\u00017\u00017\u00037\u02d6"+
		"\b7\u00018\u00018\u00019\u00019\u00039\u02dc\b9\u0001:\u0001:\u0001:\u0001"+
		":\u0003:\u02e2\b:\u0003:\u02e4\b:\u0001;\u0001;\u0001;\u0001;\u0001;\u0001"+
		";\u0001;\u0001;\u0001;\u0003;\u02ef\b;\u0001;\u0001;\u0001;\u0003;\u02f4"+
		"\b;\u0001<\u0001<\u0001<\u0005<\u02f9\b<\n<\f<\u02fc\t<\u0001<\u0003<"+
		"\u02ff\b<\u0001=\u0001=\u0001=\u0005=\u0304\b=\n=\f=\u0307\t=\u0001=\u0003"+
		"=\u030a\b=\u0001>\u0001>\u0001>\u0005>\u030f\b>\n>\f>\u0312\t>\u0001>"+
		"\u0003>\u0315\b>\u0001?\u0001?\u0001?\u0005?\u031a\b?\n?\f?\u031d\t?\u0001"+
		"?\u0003?\u0320\b?\u0001@\u0001@\u0001@\u0005@\u0325\b@\n@\f@\u0328\t@"+
		"\u0001@\u0003@\u032b\b@\u0001A\u0001A\u0001A\u0005A\u0330\bA\nA\fA\u0333"+
		"\tA\u0001A\u0003A\u0336\bA\u0001B\u0001B\u0001B\u0005B\u033b\bB\nB\fB"+
		"\u033e\tB\u0001B\u0003B\u0341\bB\u0001B\u0000\u0001LC\u0000\u0002\u0004"+
		"\u0006\b\n\f\u000e\u0010\u0012\u0014\u0016\u0018\u001a\u001c\u001e \""+
		"$&(*,.02468:<>@BDFHJLNPRTVXZ\\^`bdfhjlnprtvxz|~\u0080\u0082\u0084\u0000"+
		"\u000b\u0001\u0000\u001a\u001b\u0002\u0000\u001f\u001f!!\u0002\u0000C"+
		"CUU\u0002\u0000GGRS\u0002\u000088Yb\u0002\u0000HHJK\u0002\u0000EEGG\u0001"+
		"\u000001\u0001\u0000\u001d\u001e\u0002\u000046^^\u0001\u0000^_\u0399\u0000"+
		"\u0087\u0001\u0000\u0000\u0000\u0002\u0091\u0001\u0000\u0000\u0000\u0004"+
		"\u00bc\u0001\u0000\u0000\u0000\u0006\u00be\u0001\u0000\u0000\u0000\b\u00c2"+
		"\u0001\u0000\u0000\u0000\n\u00cb\u0001\u0000\u0000\u0000\f\u00d0\u0001"+
		"\u0000\u0000\u0000\u000e\u00d2\u0001\u0000\u0000\u0000\u0010\u00d6\u0001"+
		"\u0000\u0000\u0000\u0012\u00da\u0001\u0000\u0000\u0000\u0014\u00dd\u0001"+
		"\u0000\u0000\u0000\u0016\u00e0\u0001\u0000\u0000\u0000\u0018\u00e3\u0001"+
		"\u0000\u0000\u0000\u001a\u00f1\u0001\u0000\u0000\u0000\u001c\u00fa\u0001"+
		"\u0000\u0000\u0000\u001e\u0101\u0001\u0000\u0000\u0000 \u0107\u0001\u0000"+
		"\u0000\u0000\"\u011a\u0001\u0000\u0000\u0000$\u011c\u0001\u0000\u0000"+
		"\u0000&\u0122\u0001\u0000\u0000\u0000(\u0128\u0001\u0000\u0000\u0000*"+
		"\u012f\u0001\u0000\u0000\u0000,\u015e\u0001\u0000\u0000\u0000.\u0160\u0001"+
		"\u0000\u0000\u00000\u0167\u0001\u0000\u0000\u00002\u016b\u0001\u0000\u0000"+
		"\u00004\u0173\u0001\u0000\u0000\u00006\u017c\u0001\u0000\u0000\u00008"+
		"\u0183\u0001\u0000\u0000\u0000:\u018b\u0001\u0000\u0000\u0000<\u0192\u0001"+
		"\u0000\u0000\u0000>\u0196\u0001\u0000\u0000\u0000@\u01a2\u0001\u0000\u0000"+
		"\u0000B\u01ae\u0001\u0000\u0000\u0000D\u01ba\u0001\u0000\u0000\u0000F"+
		"\u01c2\u0001\u0000\u0000\u0000H\u01c5\u0001\u0000\u0000\u0000J\u01cc\u0001"+
		"\u0000\u0000\u0000L\u01fa\u0001\u0000\u0000\u0000N\u0224\u0001\u0000\u0000"+
		"\u0000P\u022f\u0001\u0000\u0000\u0000R\u0231\u0001\u0000\u0000\u0000T"+
		"\u0235\u0001\u0000\u0000\u0000V\u023f\u0001\u0000\u0000\u0000X\u024d\u0001"+
		"\u0000\u0000\u0000Z\u0263\u0001\u0000\u0000\u0000\\\u027a\u0001\u0000"+
		"\u0000\u0000^\u0281\u0001\u0000\u0000\u0000`\u0291\u0001\u0000\u0000\u0000"+
		"b\u02b3\u0001\u0000\u0000\u0000d\u02b5\u0001\u0000\u0000\u0000f\u02b9"+
		"\u0001\u0000\u0000\u0000h\u02c0\u0001\u0000\u0000\u0000j\u02cd\u0001\u0000"+
		"\u0000\u0000l\u02d1\u0001\u0000\u0000\u0000n\u02d5\u0001\u0000\u0000\u0000"+
		"p\u02d7\u0001\u0000\u0000\u0000r\u02db\u0001\u0000\u0000\u0000t\u02e3"+
		"\u0001\u0000\u0000\u0000v\u02f3\u0001\u0000\u0000\u0000x\u02f5\u0001\u0000"+
		"\u0000\u0000z\u0300\u0001\u0000\u0000\u0000|\u030b\u0001\u0000\u0000\u0000"+
		"~\u0316\u0001\u0000\u0000\u0000\u0080\u0321\u0001\u0000\u0000\u0000\u0082"+
		"\u032c\u0001\u0000\u0000\u0000\u0084\u0337\u0001\u0000\u0000\u0000\u0086"+
		"\u0088\u0003\u0002\u0001\u0000\u0087\u0086\u0001\u0000\u0000\u0000\u0087"+
		"\u0088\u0001\u0000\u0000\u0000\u0088\u008c\u0001\u0000\u0000\u0000\u0089"+
		"\u008b\u0003\f\u0006\u0000\u008a\u0089\u0001\u0000\u0000\u0000\u008b\u008e"+
		"\u0001\u0000\u0000\u0000\u008c\u008a\u0001\u0000\u0000\u0000\u008c\u008d"+
		"\u0001\u0000\u0000\u0000\u008d\u008f\u0001\u0000\u0000\u0000\u008e\u008c"+
		"\u0001\u0000\u0000\u0000\u008f\u0090\u0005\u0000\u0000\u0001\u0090\u0001"+
		"\u0001\u0000\u0000\u0000\u0091\u0092\u0005\u0001\u0000\u0000\u0092\u0093"+
		"\u0005h\u0000\u0000\u0093\u0094\u0005@\u0000\u0000\u0094\u0003\u0001\u0000"+
		"\u0000\u0000\u0095\u00bd\u0003\n\u0005\u0000\u0096\u0098\u0003\u0006\u0003"+
		"\u0000\u0097\u0096\u0001\u0000\u0000\u0000\u0098\u009b\u0001\u0000\u0000"+
		"\u0000\u0099\u0097\u0001\u0000\u0000\u0000\u0099\u009a\u0001\u0000\u0000"+
		"\u0000\u009a\u00ba\u0001\u0000\u0000\u0000\u009b\u0099\u0001\u0000\u0000"+
		"\u0000\u009c\u00bb\u00032\u0019\u0000\u009d\u00bb\u0003D\"\u0000\u009e"+
		"\u00bb\u0003$\u0012\u0000\u009f\u00bb\u0003&\u0013\u0000\u00a0\u00bb\u0003"+
		"\u0012\t\u0000\u00a1\u00bb\u0003H$\u0000\u00a2\u00bb\u0003\u000e\u0007"+
		"\u0000\u00a3\u00bb\u00034\u001a\u0000\u00a4\u00bb\u00036\u001b\u0000\u00a5"+
		"\u00bb\u0003\u0014\n\u0000\u00a6\u00bb\u0003>\u001f\u0000\u00a7\u00bb"+
		"\u0003J%\u0000\u00a8\u00bb\u0003(\u0014\u0000\u00a9\u00bb\u0003\u0016"+
		"\u000b\u0000\u00aa\u00bb\u0003F#\u0000\u00ab\u00bb\u0003@ \u0000\u00ac"+
		"\u00bb\u0003\u0018\f\u0000\u00ad\u00bb\u0003,\u0016\u0000\u00ae\u00bb"+
		"\u0003B!\u0000\u00af\u00bb\u0003\u001a\r\u0000\u00b0\u00bb\u0003\u0010"+
		"\b\u0000\u00b1\u00bb\u00038\u001c\u0000\u00b2\u00bb\u0003.\u0017\u0000"+
		"\u00b3\u00bb\u0003*\u0015\u0000\u00b4\u00bb\u0003:\u001d\u0000\u00b5\u00bb"+
		"\u0003<\u001e\u0000\u00b6\u00bb\u00030\u0018\u0000\u00b7\u00bb\u0003\u001c"+
		"\u000e\u0000\u00b8\u00bb\u0003 \u0010\u0000\u00b9\u00bb\u0003\u001e\u000f"+
		"\u0000\u00ba\u009c\u0001\u0000\u0000\u0000\u00ba\u009d\u0001\u0000\u0000"+
		"\u0000\u00ba\u009e\u0001\u0000\u0000\u0000\u00ba\u009f\u0001\u0000\u0000"+
		"\u0000\u00ba\u00a0\u0001\u0000\u0000\u0000\u00ba\u00a1\u0001\u0000\u0000"+
		"\u0000\u00ba\u00a2\u0001\u0000\u0000\u0000\u00ba\u00a3\u0001\u0000\u0000"+
		"\u0000\u00ba\u00a4\u0001\u0000\u0000\u0000\u00ba\u00a5\u0001\u0000\u0000"+
		"\u0000\u00ba\u00a6\u0001\u0000\u0000\u0000\u00ba\u00a7\u0001\u0000\u0000"+
		"\u0000\u00ba\u00a8\u0001\u0000\u0000\u0000\u00ba\u00a9\u0001\u0000\u0000"+
		"\u0000\u00ba\u00aa\u0001\u0000\u0000\u0000\u00ba\u00ab\u0001\u0000\u0000"+
		"\u0000\u00ba\u00ac\u0001\u0000\u0000\u0000\u00ba\u00ad\u0001\u0000\u0000"+
		"\u0000\u00ba\u00ae\u0001\u0000\u0000\u0000\u00ba\u00af\u0001\u0000\u0000"+
		"\u0000\u00ba\u00b0\u0001\u0000\u0000\u0000\u00ba\u00b1\u0001\u0000\u0000"+
		"\u0000\u00ba\u00b2\u0001\u0000\u0000\u0000\u00ba\u00b3\u0001\u0000\u0000"+
		"\u0000\u00ba\u00b4\u0001\u0000\u0000\u0000\u00ba\u00b5\u0001\u0000\u0000"+
		"\u0000\u00ba\u00b6\u0001\u0000\u0000\u0000\u00ba\u00b7\u0001\u0000\u0000"+
		"\u0000\u00ba\u00b8\u0001\u0000\u0000\u0000\u00ba\u00b9\u0001\u0000\u0000"+
		"\u0000\u00bb\u00bd\u0001\u0000\u0000\u0000\u00bc\u0095\u0001\u0000\u0000"+
		"\u0000\u00bc\u0099\u0001\u0000\u0000\u0000\u00bd\u0005\u0001\u0000\u0000"+
		"\u0000\u00be\u00c0\u0005\u0019\u0000\u0000\u00bf\u00c1\u0005m\u0000\u0000"+
		"\u00c0\u00bf\u0001\u0000\u0000\u0000\u00c0\u00c1\u0001\u0000\u0000\u0000"+
		"\u00c1\u0007\u0001\u0000\u0000\u0000\u00c2\u00c6\u0005;\u0000\u0000\u00c3"+
		"\u00c5\u0003\f\u0006\u0000\u00c4\u00c3\u0001\u0000\u0000\u0000\u00c5\u00c8"+
		"\u0001\u0000\u0000\u0000\u00c6\u00c4\u0001\u0000\u0000\u0000\u00c6\u00c7"+
		"\u0001\u0000\u0000\u0000\u00c7\u00c9\u0001\u0000\u0000\u0000\u00c8\u00c6"+
		"\u0001\u0000\u0000\u0000\u00c9\u00ca\u0005<\u0000\u0000\u00ca\t\u0001"+
		"\u0000\u0000\u0000\u00cb\u00cc\u0005\u0018\u0000\u0000\u00cc\u00cd\u0005"+
		"m\u0000\u0000\u00cd\u000b\u0001\u0000\u0000\u0000\u00ce\u00d1\u0003\u0004"+
		"\u0002\u0000\u00cf\u00d1\u0003\b\u0004\u0000\u00d0\u00ce\u0001\u0000\u0000"+
		"\u0000\u00d0\u00cf\u0001\u0000\u0000\u0000\u00d1\r\u0001\u0000\u0000\u0000"+
		"\u00d2\u00d3\u0005\u0003\u0000\u0000\u00d3\u00d4\u0005j\u0000\u0000\u00d4"+
		"\u00d5\u0005@\u0000\u0000\u00d5\u000f\u0001\u0000\u0000\u0000\u00d6\u00d7"+
		"\u0005\u0002\u0000\u0000\u00d7\u00d8\u0005j\u0000\u0000\u00d8\u00d9\u0005"+
		"@\u0000\u0000\u00d9\u0011\u0001\u0000\u0000\u0000\u00da\u00db\u0005\u000b"+
		"\u0000\u0000\u00db\u00dc\u0005@\u0000\u0000\u00dc\u0013\u0001\u0000\u0000"+
		"\u0000\u00dd\u00de\u0005\f\u0000\u0000\u00de\u00df\u0005@\u0000\u0000"+
		"\u00df\u0015\u0001\u0000\u0000\u0000\u00e0\u00e1\u0005\u000f\u0000\u0000"+
		"\u00e1\u00e2\u0005@\u0000\u0000\u00e2\u0017\u0001\u0000\u0000\u0000\u00e3"+
		"\u00e4\u0005\u0011\u0000\u0000\u00e4\u00e5\u0003b1\u0000\u00e5\u00e6\u0005"+
		"^\u0000\u0000\u00e6\u00ed\u0005\u0013\u0000\u0000\u00e7\u00ee\u0003V+"+
		"\u0000\u00e8\u00e9\u00059\u0000\u0000\u00e9\u00ea\u0003T*\u0000\u00ea"+
		"\u00eb\u0005:\u0000\u0000\u00eb\u00ee\u0001\u0000\u0000\u0000\u00ec\u00ee"+
		"\u0003L&\u0000\u00ed\u00e7\u0001\u0000\u0000\u0000\u00ed\u00e8\u0001\u0000"+
		"\u0000\u0000\u00ed\u00ec\u0001\u0000\u0000\u0000\u00ee\u00ef\u0001\u0000"+
		"\u0000\u0000\u00ef\u00f0\u0003\f\u0006\u0000\u00f0\u0019\u0001\u0000\u0000"+
		"\u0000\u00f1\u00f2\u0005\r\u0000\u0000\u00f2\u00f3\u0005=\u0000\u0000"+
		"\u00f3\u00f4\u0003L&\u0000\u00f4\u00f5\u0005>\u0000\u0000\u00f5\u00f8"+
		"\u0003\f\u0006\u0000\u00f6\u00f7\u0005\u000e\u0000\u0000\u00f7\u00f9\u0003"+
		"\f\u0006\u0000\u00f8\u00f6\u0001\u0000\u0000\u0000\u00f8\u00f9\u0001\u0000"+
		"\u0000\u0000\u00f9\u001b\u0001\u0000\u0000\u0000\u00fa\u00fd\u0005\u0010"+
		"\u0000\u0000\u00fb\u00fe\u0003L&\u0000\u00fc\u00fe\u0003R)\u0000\u00fd"+
		"\u00fb\u0001\u0000\u0000\u0000\u00fd\u00fc\u0001\u0000\u0000\u0000\u00fd"+
		"\u00fe\u0001\u0000\u0000\u0000\u00fe\u00ff\u0001\u0000\u0000\u0000\u00ff"+
		"\u0100\u0005@\u0000\u0000\u0100\u001d\u0001\u0000\u0000\u0000\u0101\u0102"+
		"\u0005\u0012\u0000\u0000\u0102\u0103\u0005=\u0000\u0000\u0103\u0104\u0003"+
		"L&\u0000\u0104\u0105\u0005>\u0000\u0000\u0105\u0106\u0003\f\u0006\u0000"+
		"\u0106\u001f\u0001\u0000\u0000\u0000\u0107\u0108\u0005\u0014\u0000\u0000"+
		"\u0108\u0109\u0005=\u0000\u0000\u0109\u010a\u0003L&\u0000\u010a\u010b"+
		"\u0005>\u0000\u0000\u010b\u010f\u0005;\u0000\u0000\u010c\u010e\u0003\""+
		"\u0011\u0000\u010d\u010c\u0001\u0000\u0000\u0000\u010e\u0111\u0001\u0000"+
		"\u0000\u0000\u010f\u010d\u0001\u0000\u0000\u0000\u010f\u0110\u0001\u0000"+
		"\u0000\u0000\u0110\u0112\u0001\u0000\u0000\u0000\u0111\u010f\u0001\u0000"+
		"\u0000\u0000\u0112\u0113\u0005<\u0000\u0000\u0113!\u0001\u0000\u0000\u0000"+
		"\u0114\u0115\u0005\u0015\u0000\u0000\u0115\u0116\u0003~?\u0000\u0116\u0117"+
		"\u0003\b\u0004\u0000\u0117\u011b\u0001\u0000\u0000\u0000\u0118\u0119\u0005"+
		"\u0016\u0000\u0000\u0119\u011b\u0003\b\u0004\u0000\u011a\u0114\u0001\u0000"+
		"\u0000\u0000\u011a\u0118\u0001\u0000\u0000\u0000\u011b#\u0001\u0000\u0000"+
		"\u0000\u011c\u011e\u00057\u0000\u0000\u011d\u011f\u0003\u0082A\u0000\u011e"+
		"\u011d\u0001\u0000\u0000\u0000\u011e\u011f\u0001\u0000\u0000\u0000\u011f"+
		"\u0120\u0001\u0000\u0000\u0000\u0120\u0121\u0005@\u0000\u0000\u0121%\u0001"+
		"\u0000\u0000\u0000\u0122\u0124\u0005\t\u0000\u0000\u0123\u0125\u0003j"+
		"5\u0000\u0124\u0123\u0001\u0000\u0000\u0000\u0124\u0125\u0001\u0000\u0000"+
		"\u0000\u0125\u0126\u0001\u0000\u0000\u0000\u0126\u0127\u0003\b\u0004\u0000"+
		"\u0127\'\u0001\u0000\u0000\u0000\u0128\u0129\u00054\u0000\u0000\u0129"+
		"\u012b\u0003j5\u0000\u012a\u012c\u0003\u0082A\u0000\u012b\u012a\u0001"+
		"\u0000\u0000\u0000\u012b\u012c\u0001\u0000\u0000\u0000\u012c\u012d\u0001"+
		"\u0000\u0000\u0000\u012d\u012e\u0005@\u0000\u0000\u012e)\u0001\u0000\u0000"+
		"\u0000\u012f\u0131\u0005\u0017\u0000\u0000\u0130\u0132\u0003\u0082A\u0000"+
		"\u0131\u0130\u0001\u0000\u0000\u0000\u0131\u0132\u0001\u0000\u0000\u0000"+
		"\u0132\u0133\u0001\u0000\u0000\u0000\u0133\u0134\u0005@\u0000\u0000\u0134"+
		"+\u0001\u0000\u0000\u0000\u0135\u0137\u0003`0\u0000\u0136\u0135\u0001"+
		"\u0000\u0000\u0000\u0137\u013a\u0001\u0000\u0000\u0000\u0138\u0136\u0001"+
		"\u0000\u0000\u0000\u0138\u0139\u0001\u0000\u0000\u0000\u0139\u013b\u0001"+
		"\u0000\u0000\u0000\u013a\u0138\u0001\u0000\u0000\u0000\u013b\u0141\u0005"+
		"^\u0000\u0000\u013c\u013e\u0005=\u0000\u0000\u013d\u013f\u0003~?\u0000"+
		"\u013e\u013d\u0001\u0000\u0000\u0000\u013e\u013f\u0001\u0000\u0000\u0000"+
		"\u013f\u0140\u0001\u0000\u0000\u0000\u0140\u0142\u0005>\u0000\u0000\u0141"+
		"\u013c\u0001\u0000\u0000\u0000\u0141\u0142\u0001\u0000\u0000\u0000\u0142"+
		"\u0144\u0001\u0000\u0000\u0000\u0143\u0145\u0003j5\u0000\u0144\u0143\u0001"+
		"\u0000\u0000\u0000\u0144\u0145\u0001\u0000\u0000\u0000\u0145\u0146\u0001"+
		"\u0000\u0000\u0000\u0146\u0147\u0003\u0082A\u0000\u0147\u0148\u0005@\u0000"+
		"\u0000\u0148\u015f\u0001\u0000\u0000\u0000\u0149\u014b\u0003`0\u0000\u014a"+
		"\u0149\u0001\u0000\u0000\u0000\u014b\u014e\u0001\u0000\u0000\u0000\u014c"+
		"\u014a\u0001\u0000\u0000\u0000\u014c\u014d\u0001\u0000\u0000\u0000\u014d"+
		"\u014f\u0001\u0000\u0000\u0000\u014e\u014c\u0001\u0000\u0000\u0000\u014f"+
		"\u0155\u0005-\u0000\u0000\u0150\u0152\u0005=\u0000\u0000\u0151\u0153\u0003"+
		"~?\u0000\u0152\u0151\u0001\u0000\u0000\u0000\u0152\u0153\u0001\u0000\u0000"+
		"\u0000\u0153\u0154\u0001\u0000\u0000\u0000\u0154\u0156\u0005>\u0000\u0000"+
		"\u0155\u0150\u0001\u0000\u0000\u0000\u0155\u0156\u0001\u0000\u0000\u0000"+
		"\u0156\u0158\u0001\u0000\u0000\u0000\u0157\u0159\u0003j5\u0000\u0158\u0157"+
		"\u0001\u0000\u0000\u0000\u0158\u0159\u0001\u0000\u0000\u0000\u0159\u015b"+
		"\u0001\u0000\u0000\u0000\u015a\u015c\u0003\u0082A\u0000\u015b\u015a\u0001"+
		"\u0000\u0000\u0000\u015b\u015c\u0001\u0000\u0000\u0000\u015c\u015d\u0001"+
		"\u0000\u0000\u0000\u015d\u015f\u0005@\u0000\u0000\u015e\u0138\u0001\u0000"+
		"\u0000\u0000\u015e\u014c\u0001\u0000\u0000\u0000\u015f-\u0001\u0000\u0000"+
		"\u0000\u0160\u0163\u0003R)\u0000\u0161\u0162\u0005D\u0000\u0000\u0162"+
		"\u0164\u0003\\.\u0000\u0163\u0161\u0001\u0000\u0000\u0000\u0163\u0164"+
		"\u0001\u0000\u0000\u0000\u0164\u0165\u0001\u0000\u0000\u0000\u0165\u0166"+
		"\u0005@\u0000\u0000\u0166/\u0001\u0000\u0000\u0000\u0167\u0168\u00055"+
		"\u0000\u0000\u0168\u0169\u0003r9\u0000\u0169\u016a\u0005@\u0000\u0000"+
		"\u016a1\u0001\u0000\u0000\u0000\u016b\u016c\u0005\n\u0000\u0000\u016c"+
		"\u016d\u0005^\u0000\u0000\u016d\u016e\u0005C\u0000\u0000\u016e\u016f\u0003"+
		"N\'\u0000\u016f\u0170\u0005@\u0000\u0000\u01703\u0001\u0000\u0000\u0000"+
		"\u0171\u0174\u0003b1\u0000\u0172\u0174\u0003f3\u0000\u0173\u0171\u0001"+
		"\u0000\u0000\u0000\u0173\u0172\u0001\u0000\u0000\u0000\u0174\u0175\u0001"+
		"\u0000\u0000\u0000\u0175\u0178\u0005^\u0000\u0000\u0176\u0177\u0005C\u0000"+
		"\u0000\u0177\u0179\u0003P(\u0000\u0178\u0176\u0001\u0000\u0000\u0000\u0178"+
		"\u0179\u0001\u0000\u0000\u0000\u0179\u017a\u0001\u0000\u0000\u0000\u017a"+
		"\u017b\u0005@\u0000\u0000\u017b5\u0001\u0000\u0000\u0000\u017c\u017d\u0005"+
		"\u001c\u0000\u0000\u017d\u017e\u0003b1\u0000\u017e\u017f\u0005^\u0000"+
		"\u0000\u017f\u0180\u0005C\u0000\u0000\u0180\u0181\u0003P(\u0000\u0181"+
		"\u0182\u0005@\u0000\u0000\u01827\u0001\u0000\u0000\u0000\u0183\u0186\u0007"+
		"\u0000\u0000\u0000\u0184\u0187\u0003b1\u0000\u0185\u0187\u0003f3\u0000"+
		"\u0186\u0184\u0001\u0000\u0000\u0000\u0186\u0185\u0001\u0000\u0000\u0000"+
		"\u0187\u0188\u0001\u0000\u0000\u0000\u0188\u0189\u0005^\u0000\u0000\u0189"+
		"\u018a\u0005@\u0000\u0000\u018a9\u0001\u0000\u0000\u0000\u018b\u018c\u0007"+
		"\u0001\u0000\u0000\u018c\u018e\u0005^\u0000\u0000\u018d\u018f\u0003j5"+
		"\u0000\u018e\u018d\u0001\u0000\u0000\u0000\u018e\u018f\u0001\u0000\u0000"+
		"\u0000\u018f\u0190\u0001\u0000\u0000\u0000\u0190\u0191\u0005@\u0000\u0000"+
		"\u0191;\u0001\u0000\u0000\u0000\u0192\u0193\u0003d2\u0000\u0193\u0194"+
		"\u0005^\u0000\u0000\u0194\u0195\u0005@\u0000\u0000\u0195=\u0001\u0000"+
		"\u0000\u0000\u0196\u0197\u0005\u0004\u0000\u0000\u0197\u0198\u0005^\u0000"+
		"\u0000\u0198\u019a\u0005=\u0000\u0000\u0199\u019b\u0003x<\u0000\u019a"+
		"\u0199\u0001\u0000\u0000\u0000\u019a\u019b\u0001\u0000\u0000\u0000\u019b"+
		"\u019c\u0001\u0000\u0000\u0000\u019c\u019e\u0005>\u0000\u0000\u019d\u019f"+
		"\u0003^/\u0000\u019e\u019d\u0001\u0000\u0000\u0000\u019e\u019f\u0001\u0000"+
		"\u0000\u0000\u019f\u01a0\u0001\u0000\u0000\u0000\u01a0\u01a1\u0003\b\u0004"+
		"\u0000\u01a1?\u0001\u0000\u0000\u0000\u01a2\u01a3\u0005\b\u0000\u0000"+
		"\u01a3\u01a4\u0005^\u0000\u0000\u01a4\u01a6\u0005=\u0000\u0000\u01a5\u01a7"+
		"\u0003\u0084B\u0000\u01a6\u01a5\u0001\u0000\u0000\u0000\u01a6\u01a7\u0001"+
		"\u0000\u0000\u0000\u01a7\u01a8\u0001\u0000\u0000\u0000\u01a8\u01aa\u0005"+
		">\u0000\u0000\u01a9\u01ab\u0003^/\u0000\u01aa\u01a9\u0001\u0000\u0000"+
		"\u0000\u01aa\u01ab\u0001\u0000\u0000\u0000\u01ab\u01ac\u0001\u0000\u0000"+
		"\u0000\u01ac\u01ad\u0005@\u0000\u0000\u01adA\u0001\u0000\u0000\u0000\u01ae"+
		"\u01af\u0005\u0007\u0000\u0000\u01af\u01b5\u0005^\u0000\u0000\u01b0\u01b2"+
		"\u0005=\u0000\u0000\u01b1\u01b3\u0003\u0080@\u0000\u01b2\u01b1\u0001\u0000"+
		"\u0000\u0000\u01b2\u01b3\u0001\u0000\u0000\u0000\u01b3\u01b4\u0001\u0000"+
		"\u0000\u0000\u01b4\u01b6\u0005>\u0000\u0000\u01b5\u01b0\u0001\u0000\u0000"+
		"\u0000\u01b5\u01b6\u0001\u0000\u0000\u0000\u01b6\u01b7\u0001\u0000\u0000"+
		"\u0000\u01b7\u01b8\u0003\u0080@\u0000\u01b8\u01b9\u0003\b\u0004\u0000"+
		"\u01b9C\u0001\u0000\u0000\u0000\u01ba\u01bb\u0003\\.\u0000\u01bb\u01be"+
		"\u0007\u0002\u0000\u0000\u01bc\u01bf\u0003L&\u0000\u01bd\u01bf\u0003R"+
		")\u0000\u01be\u01bc\u0001\u0000\u0000\u0000\u01be\u01bd\u0001\u0000\u0000"+
		"\u0000\u01bf\u01c0\u0001\u0000\u0000\u0000\u01c0\u01c1\u0005@\u0000\u0000"+
		"\u01c1E\u0001\u0000\u0000\u0000\u01c2\u01c3\u0003L&\u0000\u01c3\u01c4"+
		"\u0005@\u0000\u0000\u01c4G\u0001\u0000\u0000\u0000\u01c5\u01c6\u0005\u0005"+
		"\u0000\u0000\u01c6\u01c8\u0005;\u0000\u0000\u01c7\u01c9\u0005r\u0000\u0000"+
		"\u01c8\u01c7\u0001\u0000\u0000\u0000\u01c8\u01c9\u0001\u0000\u0000\u0000"+
		"\u01c9\u01ca\u0001\u0000\u0000\u0000\u01ca\u01cb\u0005<\u0000\u0000\u01cb"+
		"I\u0001\u0000\u0000\u0000\u01cc\u01cd\u0005\u0006\u0000\u0000\u01cd\u01d3"+
		"\u0003l6\u0000\u01ce\u01d0\u0005=\u0000\u0000\u01cf\u01d1\u0003z=\u0000"+
		"\u01d0\u01cf\u0001\u0000\u0000\u0000\u01d0\u01d1\u0001\u0000\u0000\u0000"+
		"\u01d1\u01d2\u0001\u0000\u0000\u0000\u01d2\u01d4\u0005>\u0000\u0000\u01d3"+
		"\u01ce\u0001\u0000\u0000\u0000\u01d3\u01d4\u0001\u0000\u0000\u0000\u01d4"+
		"\u01d5\u0001\u0000\u0000\u0000\u01d5\u01d7\u0003|>\u0000\u01d6\u01d8\u0003"+
		"^/\u0000\u01d7\u01d6\u0001\u0000\u0000\u0000\u01d7\u01d8\u0001\u0000\u0000"+
		"\u0000\u01d8\u01d9\u0001\u0000\u0000\u0000\u01d9\u01db\u0005;\u0000\u0000"+
		"\u01da\u01dc\u0005r\u0000\u0000\u01db\u01da\u0001\u0000\u0000\u0000\u01db"+
		"\u01dc\u0001\u0000\u0000\u0000\u01dc\u01dd\u0001\u0000\u0000\u0000\u01dd"+
		"\u01de\u0005<\u0000\u0000\u01deK\u0001\u0000\u0000\u0000\u01df\u01e0\u0006"+
		"&\uffff\uffff\u0000\u01e0\u01e1\u0005=\u0000\u0000\u01e1\u01e2\u0003L"+
		"&\u0000\u01e2\u01e3\u0005>\u0000\u0000\u01e3\u01fb\u0001\u0000\u0000\u0000"+
		"\u01e4\u01e5\u0007\u0003\u0000\u0000\u01e5\u01fb\u0003L&\u000f\u01e6\u01e9"+
		"\u0003b1\u0000\u01e7\u01e9\u0003f3\u0000\u01e8\u01e6\u0001\u0000\u0000"+
		"\u0000\u01e8\u01e7\u0001\u0000\u0000\u0000\u01e9\u01ea\u0001\u0000\u0000"+
		"\u0000\u01ea\u01eb\u0005=\u0000\u0000\u01eb\u01ec\u0003L&\u0000\u01ec"+
		"\u01ed\u0005>\u0000\u0000\u01ed\u01fb\u0001\u0000\u0000\u0000\u01ee\u01ef"+
		"\u00053\u0000\u0000\u01ef\u01f0\u0005=\u0000\u0000\u01f0\u01f1\u0003\b"+
		"\u0004\u0000\u01f1\u01f2\u0005>\u0000\u0000\u01f2\u01fb\u0001\u0000\u0000"+
		"\u0000\u01f3\u01f4\u0005^\u0000\u0000\u01f4\u01f6\u0005=\u0000\u0000\u01f5"+
		"\u01f7\u0003~?\u0000\u01f6\u01f5\u0001\u0000\u0000\u0000\u01f6\u01f7\u0001"+
		"\u0000\u0000\u0000\u01f7\u01f8\u0001\u0000\u0000\u0000\u01f8\u01fb\u0005"+
		">\u0000\u0000\u01f9\u01fb\u0007\u0004\u0000\u0000\u01fa\u01df\u0001\u0000"+
		"\u0000\u0000\u01fa\u01e4\u0001\u0000\u0000\u0000\u01fa\u01e8\u0001\u0000"+
		"\u0000\u0000\u01fa\u01ee\u0001\u0000\u0000\u0000\u01fa\u01f3\u0001\u0000"+
		"\u0000\u0000\u01fa\u01f9\u0001\u0000\u0000\u0000\u01fb\u0221\u0001\u0000"+
		"\u0000\u0000\u01fc\u01fd\n\u0010\u0000\u0000\u01fd\u01fe\u0005I\u0000"+
		"\u0000\u01fe\u0220\u0003L&\u0010\u01ff\u0200\n\u000e\u0000\u0000\u0200"+
		"\u0201\u0007\u0005\u0000\u0000\u0201\u0220\u0003L&\u000f\u0202\u0203\n"+
		"\r\u0000\u0000\u0203\u0204\u0007\u0006\u0000\u0000\u0204\u0220\u0003L"+
		"&\u000e\u0205\u0206\n\f\u0000\u0000\u0206\u0207\u0005W\u0000\u0000\u0207"+
		"\u0220\u0003L&\r\u0208\u0209\n\u000b\u0000\u0000\u0209\u020a\u0005V\u0000"+
		"\u0000\u020a\u0220\u0003L&\f\u020b\u020c\n\n\u0000\u0000\u020c\u020d\u0005"+
		"T\u0000\u0000\u020d\u0220\u0003L&\u000b\u020e\u020f\n\t\u0000\u0000\u020f"+
		"\u0210\u0005N\u0000\u0000\u0210\u0220\u0003L&\n\u0211\u0212\n\b\u0000"+
		"\u0000\u0212\u0213\u0005P\u0000\u0000\u0213\u0220\u0003L&\t\u0214\u0215"+
		"\n\u0007\u0000\u0000\u0215\u0216\u0005L\u0000\u0000\u0216\u0220\u0003"+
		"L&\b\u0217\u0218\n\u0006\u0000\u0000\u0218\u0219\u0005O\u0000\u0000\u0219"+
		"\u0220\u0003L&\u0007\u021a\u021b\n\u0005\u0000\u0000\u021b\u021c\u0005"+
		"M\u0000\u0000\u021c\u0220\u0003L&\u0006\u021d\u021e\n\u0011\u0000\u0000"+
		"\u021e\u0220\u0003Z-\u0000\u021f\u01fc\u0001\u0000\u0000\u0000\u021f\u01ff"+
		"\u0001\u0000\u0000\u0000\u021f\u0202\u0001\u0000\u0000\u0000\u021f\u0205"+
		"\u0001\u0000\u0000\u0000\u021f\u0208\u0001\u0000\u0000\u0000\u021f\u020b"+
		"\u0001\u0000\u0000\u0000\u021f\u020e\u0001\u0000\u0000\u0000\u021f\u0211"+
		"\u0001\u0000\u0000\u0000\u021f\u0214\u0001\u0000\u0000\u0000\u021f\u0217"+
		"\u0001\u0000\u0000\u0000\u021f\u021a\u0001\u0000\u0000\u0000\u021f\u021d"+
		"\u0001\u0000\u0000\u0000\u0220\u0223\u0001\u0000\u0000\u0000\u0221\u021f"+
		"\u0001\u0000\u0000\u0000\u0221\u0222\u0001\u0000\u0000\u0000\u0222M\u0001"+
		"\u0000\u0000\u0000\u0223\u0221\u0001\u0000\u0000\u0000\u0224\u0229\u0003"+
		"L&\u0000\u0225\u0226\u0005F\u0000\u0000\u0226\u0228\u0003L&\u0000\u0227"+
		"\u0225\u0001\u0000\u0000\u0000\u0228\u022b\u0001\u0000\u0000\u0000\u0229"+
		"\u0227\u0001\u0000\u0000\u0000\u0229\u022a\u0001\u0000\u0000\u0000\u022a"+
		"O\u0001\u0000\u0000\u0000\u022b\u0229\u0001\u0000\u0000\u0000\u022c\u0230"+
		"\u0003X,\u0000\u022d\u0230\u0003L&\u0000\u022e\u0230\u0003R)\u0000\u022f"+
		"\u022c\u0001\u0000\u0000\u0000\u022f\u022d\u0001\u0000\u0000\u0000\u022f"+
		"\u022e\u0001\u0000\u0000\u0000\u0230Q\u0001\u0000\u0000\u0000\u0231\u0232"+
		"\u00056\u0000\u0000\u0232\u0233\u0003r9\u0000\u0233S\u0001\u0000\u0000"+
		"\u0000\u0234\u0236\u0003L&\u0000\u0235\u0234\u0001\u0000\u0000\u0000\u0235"+
		"\u0236\u0001\u0000\u0000\u0000\u0236\u0237\u0001\u0000\u0000\u0000\u0237"+
		"\u0239\u0005?\u0000\u0000\u0238\u023a\u0003L&\u0000\u0239\u0238\u0001"+
		"\u0000\u0000\u0000\u0239\u023a\u0001\u0000\u0000\u0000\u023a\u023d\u0001"+
		"\u0000\u0000\u0000\u023b\u023c\u0005?\u0000\u0000\u023c\u023e\u0003L&"+
		"\u0000\u023d\u023b\u0001\u0000\u0000\u0000\u023d\u023e\u0001\u0000\u0000"+
		"\u0000\u023eU\u0001\u0000\u0000\u0000\u023f\u0240\u0005;\u0000\u0000\u0240"+
		"\u0245\u0003L&\u0000\u0241\u0242\u0005B\u0000\u0000\u0242\u0244\u0003"+
		"L&\u0000\u0243\u0241\u0001\u0000\u0000\u0000\u0244\u0247\u0001\u0000\u0000"+
		"\u0000\u0245\u0243\u0001\u0000\u0000\u0000\u0245\u0246\u0001\u0000\u0000"+
		"\u0000\u0246\u0249\u0001\u0000\u0000\u0000\u0247\u0245\u0001\u0000\u0000"+
		"\u0000\u0248\u024a\u0005B\u0000\u0000\u0249\u0248\u0001\u0000\u0000\u0000"+
		"\u0249\u024a\u0001\u0000\u0000\u0000\u024a\u024b\u0001\u0000\u0000\u0000"+
		"\u024b\u024c\u0005<\u0000\u0000\u024cW\u0001\u0000\u0000\u0000\u024d\u025f"+
		"\u0005;\u0000\u0000\u024e\u0251\u0003L&\u0000\u024f\u0251\u0003X,\u0000"+
		"\u0250\u024e\u0001\u0000\u0000\u0000\u0250\u024f\u0001\u0000\u0000\u0000"+
		"\u0251\u0259\u0001\u0000\u0000\u0000\u0252\u0255\u0005B\u0000\u0000\u0253"+
		"\u0256\u0003L&\u0000\u0254\u0256\u0003X,\u0000\u0255\u0253\u0001\u0000"+
		"\u0000\u0000\u0255\u0254\u0001\u0000\u0000\u0000\u0256\u0258\u0001\u0000"+
		"\u0000\u0000\u0257\u0252\u0001\u0000\u0000\u0000\u0258\u025b\u0001\u0000"+
		"\u0000\u0000\u0259\u0257\u0001\u0000\u0000\u0000\u0259\u025a\u0001\u0000"+
		"\u0000\u0000\u025a\u025d\u0001\u0000\u0000\u0000\u025b\u0259\u0001\u0000"+
		"\u0000\u0000\u025c\u025e\u0005B\u0000\u0000\u025d\u025c\u0001\u0000\u0000"+
		"\u0000\u025d\u025e\u0001\u0000\u0000\u0000\u025e\u0260\u0001\u0000\u0000"+
		"\u0000\u025f\u0250\u0001\u0000\u0000\u0000\u025f\u0260\u0001\u0000\u0000"+
		"\u0000\u0260\u0261\u0001\u0000\u0000\u0000\u0261\u0262\u0005<\u0000\u0000"+
		"\u0262Y\u0001\u0000\u0000\u0000\u0263\u0276\u00059\u0000\u0000\u0264\u0277"+
		"\u0003V+\u0000\u0265\u0268\u0003L&\u0000\u0266\u0268\u0003T*\u0000\u0267"+
		"\u0265\u0001\u0000\u0000\u0000\u0267\u0266\u0001\u0000\u0000\u0000\u0268"+
		"\u0270\u0001\u0000\u0000\u0000\u0269\u026c\u0005B\u0000\u0000\u026a\u026d"+
		"\u0003L&\u0000\u026b\u026d\u0003T*\u0000\u026c\u026a\u0001\u0000\u0000"+
		"\u0000\u026c\u026b\u0001\u0000\u0000\u0000\u026d\u026f\u0001\u0000\u0000"+
		"\u0000\u026e\u0269\u0001\u0000\u0000\u0000\u026f\u0272\u0001\u0000\u0000"+
		"\u0000\u0270\u026e\u0001\u0000\u0000\u0000\u0270\u0271\u0001\u0000\u0000"+
		"\u0000\u0271\u0274\u0001\u0000\u0000\u0000\u0272\u0270\u0001\u0000\u0000"+
		"\u0000\u0273\u0275\u0005B\u0000\u0000\u0274\u0273\u0001\u0000\u0000\u0000"+
		"\u0274\u0275\u0001\u0000\u0000\u0000\u0275\u0277\u0001\u0000\u0000\u0000"+
		"\u0276\u0264\u0001\u0000\u0000\u0000\u0276\u0267\u0001\u0000\u0000\u0000"+
		"\u0277\u0278\u0001\u0000\u0000\u0000\u0278\u0279\u0005:\u0000\u0000\u0279"+
		"[\u0001\u0000\u0000\u0000\u027a\u027e\u0005^\u0000\u0000\u027b\u027d\u0003"+
		"Z-\u0000\u027c\u027b\u0001\u0000\u0000\u0000\u027d\u0280\u0001\u0000\u0000"+
		"\u0000\u027e\u027c\u0001\u0000\u0000\u0000\u027e\u027f\u0001\u0000\u0000"+
		"\u0000\u027f]\u0001\u0000\u0000\u0000\u0280\u027e\u0001\u0000\u0000\u0000"+
		"\u0281\u0282\u0005D\u0000\u0000\u0282\u0283\u0003b1\u0000\u0283_\u0001"+
		"\u0000\u0000\u0000\u0284\u0292\u0005.\u0000\u0000\u0285\u0286\u0005/\u0000"+
		"\u0000\u0286\u0287\u0005=\u0000\u0000\u0287\u0288\u0003L&\u0000\u0288"+
		"\u0289\u0005>\u0000\u0000\u0289\u0292\u0001\u0000\u0000\u0000\u028a\u028f"+
		"\u0007\u0007\u0000\u0000\u028b\u028c\u0005=\u0000\u0000\u028c\u028d\u0003"+
		"L&\u0000\u028d\u028e\u0005>\u0000\u0000\u028e\u0290\u0001\u0000\u0000"+
		"\u0000\u028f\u028b\u0001\u0000\u0000\u0000\u028f\u0290\u0001\u0000\u0000"+
		"\u0000\u0290\u0292\u0001\u0000\u0000\u0000\u0291\u0284\u0001\u0000\u0000"+
		"\u0000\u0291\u0285\u0001\u0000\u0000\u0000\u0291\u028a\u0001\u0000\u0000"+
		"\u0000\u0292\u0293\u0001\u0000\u0000\u0000\u0293\u0294\u0005Q\u0000\u0000"+
		"\u0294a\u0001\u0000\u0000\u0000\u0295\u0297\u0005#\u0000\u0000\u0296\u0298"+
		"\u0003j5\u0000\u0297\u0296\u0001\u0000\u0000\u0000\u0297\u0298\u0001\u0000"+
		"\u0000\u0000\u0298\u02b4\u0001\u0000\u0000\u0000\u0299\u029b\u0005$\u0000"+
		"\u0000\u029a\u029c\u0003j5\u0000\u029b\u029a\u0001\u0000\u0000\u0000\u029b"+
		"\u029c\u0001\u0000\u0000\u0000\u029c\u02b4\u0001\u0000\u0000\u0000\u029d"+
		"\u029f\u0005%\u0000\u0000\u029e\u02a0\u0003j5\u0000\u029f\u029e\u0001"+
		"\u0000\u0000\u0000\u029f\u02a0\u0001\u0000\u0000\u0000\u02a0\u02b4\u0001"+
		"\u0000\u0000\u0000\u02a1\u02a3\u0005&\u0000\u0000\u02a2\u02a4\u0003j5"+
		"\u0000\u02a3\u02a2\u0001\u0000\u0000\u0000\u02a3\u02a4\u0001\u0000\u0000"+
		"\u0000\u02a4\u02b4\u0001\u0000\u0000\u0000\u02a5\u02a7\u0005\'\u0000\u0000"+
		"\u02a6\u02a8\u0003j5\u0000\u02a7\u02a6\u0001\u0000\u0000\u0000\u02a7\u02a8"+
		"\u0001\u0000\u0000\u0000\u02a8\u02b4\u0001\u0000\u0000\u0000\u02a9\u02b4"+
		"\u0005\"\u0000\u0000\u02aa\u02b4\u0005+\u0000\u0000\u02ab\u02b4\u0005"+
		",\u0000\u0000\u02ac\u02b1\u0005(\u0000\u0000\u02ad\u02ae\u00059\u0000"+
		"\u0000\u02ae\u02af\u0003b1\u0000\u02af\u02b0\u0005:\u0000\u0000\u02b0"+
		"\u02b2\u0001\u0000\u0000\u0000\u02b1\u02ad\u0001\u0000\u0000\u0000\u02b1"+
		"\u02b2\u0001\u0000\u0000\u0000\u02b2\u02b4\u0001\u0000\u0000\u0000\u02b3"+
		"\u0295\u0001\u0000\u0000\u0000\u02b3\u0299\u0001\u0000\u0000\u0000\u02b3"+
		"\u029d\u0001\u0000\u0000\u0000\u02b3\u02a1\u0001\u0000\u0000\u0000\u02b3"+
		"\u02a5\u0001\u0000\u0000\u0000\u02b3\u02a9\u0001\u0000\u0000\u0000\u02b3"+
		"\u02aa\u0001\u0000\u0000\u0000\u02b3\u02ab\u0001\u0000\u0000\u0000\u02b3"+
		"\u02ac\u0001\u0000\u0000\u0000\u02b4c\u0001\u0000\u0000\u0000\u02b5\u02b7"+
		"\u0005 \u0000\u0000\u02b6\u02b8\u0003j5\u0000\u02b7\u02b6\u0001\u0000"+
		"\u0000\u0000\u02b7\u02b8\u0001\u0000\u0000\u0000\u02b8e\u0001\u0000\u0000"+
		"\u0000\u02b9\u02ba\u0005)\u0000\u0000\u02ba\u02bb\u00059\u0000\u0000\u02bb"+
		"\u02bc\u0003b1\u0000\u02bc\u02bd\u0005B\u0000\u0000\u02bd\u02be\u0003"+
		"~?\u0000\u02be\u02bf\u0005:\u0000\u0000\u02bfg\u0001\u0000\u0000\u0000"+
		"\u02c0\u02c1\u0007\b\u0000\u0000\u02c1\u02c2\u0005)\u0000\u0000\u02c2"+
		"\u02c3\u00059\u0000\u0000\u02c3\u02c4\u0003b1\u0000\u02c4\u02c9\u0005"+
		"B\u0000\u0000\u02c5\u02ca\u0003~?\u0000\u02c6\u02c7\u00052\u0000\u0000"+
		"\u02c7\u02c8\u0005C\u0000\u0000\u02c8\u02ca\u0003L&\u0000\u02c9\u02c5"+
		"\u0001\u0000\u0000\u0000\u02c9\u02c6\u0001\u0000\u0000\u0000\u02ca\u02cb"+
		"\u0001\u0000\u0000\u0000\u02cb\u02cc\u0005:\u0000\u0000\u02cci\u0001\u0000"+
		"\u0000\u0000\u02cd\u02ce\u00059\u0000\u0000\u02ce\u02cf\u0003L&\u0000"+
		"\u02cf\u02d0\u0005:\u0000\u0000\u02d0k\u0001\u0000\u0000\u0000\u02d1\u02d2"+
		"\u0007\t\u0000\u0000\u02d2m\u0001\u0000\u0000\u0000\u02d3\u02d6\u0003"+
		"L&\u0000\u02d4\u02d6\u0003v;\u0000\u02d5\u02d3\u0001\u0000\u0000\u0000"+
		"\u02d5\u02d4\u0001\u0000\u0000\u0000\u02d6o\u0001\u0000\u0000\u0000\u02d7"+
		"\u02d8\u0007\n\u0000\u0000\u02d8q\u0001\u0000\u0000\u0000\u02d9\u02dc"+
		"\u0003\\.\u0000\u02da\u02dc\u0005_\u0000\u0000\u02db\u02d9\u0001\u0000"+
		"\u0000\u0000\u02db\u02da\u0001\u0000\u0000\u0000\u02dcs\u0001\u0000\u0000"+
		"\u0000\u02dd\u02e4\u0003b1\u0000\u02de\u02e4\u0003h4\u0000\u02df\u02e1"+
		"\u0005!\u0000\u0000\u02e0\u02e2\u0003j5\u0000\u02e1\u02e0\u0001\u0000"+
		"\u0000\u0000\u02e1\u02e2\u0001\u0000\u0000\u0000\u02e2\u02e4\u0001\u0000"+
		"\u0000\u0000\u02e3\u02dd\u0001\u0000\u0000\u0000\u02e3\u02de\u0001\u0000"+
		"\u0000\u0000\u02e3\u02df\u0001\u0000\u0000\u0000\u02e4u\u0001\u0000\u0000"+
		"\u0000\u02e5\u02e6\u0003b1\u0000\u02e6\u02e7\u0005^\u0000\u0000\u02e7"+
		"\u02f4\u0001\u0000\u0000\u0000\u02e8\u02e9\u0003d2\u0000\u02e9\u02ea\u0005"+
		"^\u0000\u0000\u02ea\u02f4\u0001\u0000\u0000\u0000\u02eb\u02ec\u0007\u0001"+
		"\u0000\u0000\u02ec\u02ee\u0005^\u0000\u0000\u02ed\u02ef\u0003j5\u0000"+
		"\u02ee\u02ed\u0001\u0000\u0000\u0000\u02ee\u02ef\u0001\u0000\u0000\u0000"+
		"\u02ef\u02f4\u0001\u0000\u0000\u0000\u02f0\u02f1\u0003h4\u0000\u02f1\u02f2"+
		"\u0005^\u0000\u0000\u02f2\u02f4\u0001\u0000\u0000\u0000\u02f3\u02e5\u0001"+
		"\u0000\u0000\u0000\u02f3\u02e8\u0001\u0000\u0000\u0000\u02f3\u02eb\u0001"+
		"\u0000\u0000\u0000\u02f3\u02f0\u0001\u0000\u0000\u0000\u02f4w\u0001\u0000"+
		"\u0000\u0000\u02f5\u02fa\u0003v;\u0000\u02f6\u02f7\u0005B\u0000\u0000"+
		"\u02f7\u02f9\u0003v;\u0000\u02f8\u02f6\u0001\u0000\u0000\u0000\u02f9\u02fc"+
		"\u0001\u0000\u0000\u0000\u02fa\u02f8\u0001\u0000\u0000\u0000\u02fa\u02fb"+
		"\u0001\u0000\u0000\u0000\u02fb\u02fe\u0001\u0000\u0000\u0000\u02fc\u02fa"+
		"\u0001\u0000\u0000\u0000\u02fd\u02ff\u0005B\u0000\u0000\u02fe\u02fd\u0001"+
		"\u0000\u0000\u0000\u02fe\u02ff\u0001\u0000\u0000\u0000\u02ffy\u0001\u0000"+
		"\u0000\u0000\u0300\u0305\u0003n7\u0000\u0301\u0302\u0005B\u0000\u0000"+
		"\u0302\u0304\u0003n7\u0000\u0303\u0301\u0001\u0000\u0000\u0000\u0304\u0307"+
		"\u0001\u0000\u0000\u0000\u0305\u0303\u0001\u0000\u0000\u0000\u0305\u0306"+
		"\u0001\u0000\u0000\u0000\u0306\u0309\u0001\u0000\u0000\u0000\u0307\u0305"+
		"\u0001\u0000\u0000\u0000\u0308\u030a\u0005B\u0000\u0000\u0309\u0308\u0001"+
		"\u0000\u0000\u0000\u0309\u030a\u0001\u0000\u0000\u0000\u030a{\u0001\u0000"+
		"\u0000\u0000\u030b\u0310\u0003p8\u0000\u030c\u030d\u0005B\u0000\u0000"+
		"\u030d\u030f\u0003p8\u0000\u030e\u030c\u0001\u0000\u0000\u0000\u030f\u0312"+
		"\u0001\u0000\u0000\u0000\u0310\u030e\u0001\u0000\u0000\u0000\u0310\u0311"+
		"\u0001\u0000\u0000\u0000\u0311\u0314\u0001\u0000\u0000\u0000\u0312\u0310"+
		"\u0001\u0000\u0000\u0000\u0313\u0315\u0005B\u0000\u0000\u0314\u0313\u0001"+
		"\u0000\u0000\u0000\u0314\u0315\u0001\u0000\u0000\u0000\u0315}\u0001\u0000"+
		"\u0000\u0000\u0316\u031b\u0003L&\u0000\u0317\u0318\u0005B\u0000\u0000"+
		"\u0318\u031a\u0003L&\u0000\u0319\u0317\u0001\u0000\u0000\u0000\u031a\u031d"+
		"\u0001\u0000\u0000\u0000\u031b\u0319\u0001\u0000\u0000\u0000\u031b\u031c"+
		"\u0001\u0000\u0000\u0000\u031c\u031f\u0001\u0000\u0000\u0000\u031d\u031b"+
		"\u0001\u0000\u0000\u0000\u031e\u0320\u0005B\u0000\u0000\u031f\u031e\u0001"+
		"\u0000\u0000\u0000\u031f\u0320\u0001\u0000\u0000\u0000\u0320\u007f\u0001"+
		"\u0000\u0000\u0000\u0321\u0326\u0005^\u0000\u0000\u0322\u0323\u0005B\u0000"+
		"\u0000\u0323\u0325\u0005^\u0000\u0000\u0324\u0322\u0001\u0000\u0000\u0000"+
		"\u0325\u0328\u0001\u0000\u0000\u0000\u0326\u0324\u0001\u0000\u0000\u0000"+
		"\u0326\u0327\u0001\u0000\u0000\u0000\u0327\u032a\u0001\u0000\u0000\u0000"+
		"\u0328\u0326\u0001\u0000\u0000\u0000\u0329\u032b\u0005B\u0000\u0000\u032a"+
		"\u0329\u0001\u0000\u0000\u0000\u032a\u032b\u0001\u0000\u0000\u0000\u032b"+
		"\u0081\u0001\u0000\u0000\u0000\u032c\u0331\u0003r9\u0000\u032d\u032e\u0005"+
		"B\u0000\u0000\u032e\u0330\u0003r9\u0000\u032f\u032d\u0001\u0000\u0000"+
		"\u0000\u0330\u0333\u0001\u0000\u0000\u0000\u0331\u032f\u0001\u0000\u0000"+
		"\u0000\u0331\u0332\u0001\u0000\u0000\u0000\u0332\u0335\u0001\u0000\u0000"+
		"\u0000\u0333\u0331\u0001\u0000\u0000\u0000\u0334\u0336\u0005B\u0000\u0000"+
		"\u0335\u0334\u0001\u0000\u0000\u0000\u0335\u0336\u0001\u0000\u0000\u0000"+
		"\u0336\u0083\u0001\u0000\u0000\u0000\u0337\u033c\u0003t:\u0000\u0338\u0339"+
		"\u0005B\u0000\u0000\u0339\u033b\u0003t:\u0000\u033a\u0338\u0001\u0000"+
		"\u0000\u0000\u033b\u033e\u0001\u0000\u0000\u0000\u033c\u033a\u0001\u0000"+
		"\u0000\u0000\u033c\u033d\u0001\u0000\u0000\u0000\u033d\u0340\u0001\u0000"+
		"\u0000\u0000\u033e\u033c\u0001\u0000\u0000\u0000\u033f\u0341\u0005B\u0000"+
		"\u0000\u0340\u033f\u0001\u0000\u0000\u0000\u0340\u0341\u0001\u0000\u0000"+
		"\u0000\u0341\u0085\u0001\u0000\u0000\u0000b\u0087\u008c\u0099\u00ba\u00bc"+
		"\u00c0\u00c6\u00d0\u00ed\u00f8\u00fd\u010f\u011a\u011e\u0124\u012b\u0131"+
		"\u0138\u013e\u0141\u0144\u014c\u0152\u0155\u0158\u015b\u015e\u0163\u0173"+
		"\u0178\u0186\u018e\u019a\u019e\u01a6\u01aa\u01b2\u01b5\u01be\u01c8\u01d0"+
		"\u01d3\u01d7\u01db\u01e8\u01f6\u01fa\u021f\u0221\u0229\u022f\u0235\u0239"+
		"\u023d\u0245\u0249\u0250\u0255\u0259\u025d\u025f\u0267\u026c\u0270\u0274"+
		"\u0276\u027e\u028f\u0291\u0297\u029b\u029f\u02a3\u02a7\u02b1\u02b3\u02b7"+
		"\u02c9\u02d5\u02db\u02e1\u02e3\u02ee\u02f3\u02fa\u02fe\u0305\u0309\u0310"+
		"\u0314\u031b\u031f\u0326\u032a\u0331\u0335\u033c\u0340";
	public static final ATN _ATN =
		new ATNDeserializer().deserialize(_serializedATN.toCharArray());
	static {
		_decisionToDFA = new DFA[_ATN.getNumberOfDecisions()];
		for (int i = 0; i < _ATN.getNumberOfDecisions(); i++) {
			_decisionToDFA[i] = new DFA(_ATN.getDecisionState(i), i);
		}
	}
}