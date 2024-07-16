program		: top_compstmt
                ;

top_compstmt	: top_stmts terms?
                ;

top_stmts	: none
                | top_stmt
                | top_stmts terms top_stmt
                ;

top_stmt	: stmt
                | keyword_BEGIN begin_block
                ;

block_open	: '{' ;

begin_block	: block_open top_compstmt '}'
                ;

bodystmt	: compstmt[body]
                  lex_ctxt[ctxt]
                  opt_rescue
                  k_else
                  compstmt[elsebody]
                  opt_ensure
                | compstmt[body]
                  lex_ctxt[ctxt]
                  opt_rescue
                  opt_ensure
                ;

compstmt	: stmts terms?
                ;

stmts		: none
                | stmt_or_begin
                | stmts terms stmt_or_begin
                ;

stmt_or_begin	: stmt
                | keyword_BEGIN
                  begin_block
                ;

allow_exits	: ;

k_END		: keyword_END lex_ctxt
                    ;

stmt		: keyword_alias fitem  fitem
                | keyword_alias tGVAR tGVAR
                | keyword_alias tGVAR tBACK_REF
                | keyword_alias tGVAR tNTH_REF
                | keyword_undef undef_list
                | stmt modifier_if expr_value
                | stmt modifier_unless expr_value
                | stmt modifier_while expr_value
                | stmt modifier_until expr_value
                | stmt modifier_rescue after_rescue stmt
                | k_END allow_exits '{' compstmt '}'
                | command_asgn
                | mlhs '=' lex_ctxt command_call
                | lhs '=' lex_ctxt mrhs
                | mlhs '=' lex_ctxt mrhs_arg modifier_rescue
                  after_rescue stmt[resbody]
                | mlhs '=' lex_ctxt mrhs_arg
                | expr
                | error
                ;

command_asgn	: lhs '=' lex_ctxt command_rhs
                | var_lhs tOP_ASGN lex_ctxt command_rhs
                | primary_value '[' opt_call_args rbracket tOP_ASGN lex_ctxt command_rhs
                | primary_value call_op ident_or_const tOP_ASGN lex_ctxt command_rhs
                | primary_value tCOLON2 tCONSTANT tOP_ASGN lex_ctxt command_rhs
                | primary_value tCOLON2 tIDENTIFIER tOP_ASGN lex_ctxt command_rhs
                | defn_head[head] f_opt_paren_args[args] '=' endless_command[bodystmt]
                | defs_head[head] f_opt_paren_args[args] '=' endless_command[bodystmt]
                | backref tOP_ASGN lex_ctxt command_rhs
                ;

endless_command : command
                | endless_command modifier_rescue after_rescue arg
                | keyword_not '\n'? endless_command
                ;

command_rhs	: command_call   %prec tOP_ASGN
                | command_call modifier_rescue after_rescue stmt
                | command_asgn
                ;

expr		: command_call
                | expr keyword_and expr
                | expr keyword_or expr
                | keyword_not '\n'? expr
                | '!' command_call
                | arg tASSOC
                  p_in_kwarg[ctxt] p_pvtbl p_pktbl
                  p_top_expr_body[body]
                | arg keyword_in
                  p_in_kwarg[ctxt] p_pvtbl p_pktbl
                  p_top_expr_body[body]
                | arg %prec tLBRACE_ARG
                ;

def_name	: fname
                ;

defn_head	: k_def def_name
                ;

defs_head	: k_def singleton dot_or_colon
                  def_name
                ;

expr_value	: expr
                | error
                ;

expr_value_do	: expr_value do
                ;

command_call	: command
                | block_command
                ;

block_command	: block_call
                | block_call call_op2 operation2 command_args
                ;

cmd_brace_block	: tLBRACE_ARG brace_body '}'
                ;

fcall		: operation
                ;

command		: fcall command_args       %prec tLOWEST
                | fcall command_args cmd_brace_block
                | primary_value call_op operation2 command_args	%prec tLOWEST
                | primary_value call_op operation2 command_args cmd_brace_block
                | primary_value tCOLON2 operation2 command_args	%prec tLOWEST
                | primary_value tCOLON2 operation2 command_args cmd_brace_block
                | primary_value tCOLON2 tCONSTANT '{' brace_body '}'
                | keyword_super command_args
                | k_yield command_args
                | k_return call_args
                | keyword_break call_args
                | keyword_next call_args
                ;

mlhs		: mlhs_basic
                | tLPAREN mlhs_inner rparen
                ;

mlhs_inner	: mlhs_basic
                | tLPAREN mlhs_inner rparen
                ;

mlhs_basic	: mlhs_head
                | mlhs_head mlhs_item
                | mlhs_head tSTAR mlhs_node
                | mlhs_head tSTAR mlhs_node ',' mlhs_post
                | mlhs_head tSTAR
                | mlhs_head tSTAR ',' mlhs_post
                | tSTAR mlhs_node
                | tSTAR mlhs_node ',' mlhs_post
                | tSTAR
                | tSTAR ',' mlhs_post
                ;

mlhs_item	: mlhs_node
                | tLPAREN mlhs_inner rparen
                ;

mlhs_head	: mlhs_item ','
                | mlhs_head mlhs_item ','
                ;

mlhs_post	: mlhs_item
                | mlhs_post ',' mlhs_item
                ;

mlhs_node	: user_variable
                | keyword_variable
                | primary_value '[' opt_call_args rbracket
                | primary_value call_op tIDENTIFIER
                | primary_value tCOLON2 tIDENTIFIER
                | primary_value call_op tCONSTANT
                | primary_value tCOLON2 tCONSTANT
                | tCOLON3 tCONSTANT
                | backref
                ;

lhs		: user_variable
                | keyword_variable
                | primary_value '[' opt_call_args rbracket
                | primary_value call_op tIDENTIFIER
                | primary_value tCOLON2 tIDENTIFIER
                | primary_value call_op tCONSTANT
                | primary_value tCOLON2 tCONSTANT
                | tCOLON3 tCONSTANT
                | backref
                ;

cname		: tIDENTIFIER
                | tCONSTANT
                ;

cpath		: tCOLON3 cname
                | cname
                | primary_value tCOLON2 cname
                ;

fname		: ident_or_const
                | tFID
                | op
                | reswords
                ;

fitem		: fname
                | symbol
                ;

undef_list	: fitem
                | undef_list ','  fitem
                ;

op		: '|'
                | '^'
                | '&'
                | tCMP
                | tEQ
                | tEQQ
                | tMATCH
                | tNMATCH
                | '>'
                | tGEQ
                | '<'
                | tLEQ
                | tNEQ
                | tLSHFT
                | tRSHFT
                | '+'
                | '-'
                | '*'
                | tSTAR
                | '/'
                | '%'
                | tPOW
                | tDSTAR
                | '!'
                | '~'
                | tUPLUS
                | tUMINUS
                | tAREF
                | tASET
                | '`'
                ;

reswords	: keyword__LINE__ | keyword__FILE__ | keyword__ENCODING__
                | keyword_BEGIN | keyword_END
                | keyword_alias | keyword_and | keyword_begin
                | keyword_break | keyword_case | keyword_class | keyword_def
                | keyword_defined | keyword_do | keyword_else | keyword_elsif
                | keyword_end | keyword_ensure | keyword_false
                | keyword_for | keyword_in | keyword_module | keyword_next
                | keyword_nil | keyword_not | keyword_or | keyword_redo
                | keyword_rescue | keyword_retry | keyword_return | keyword_self
                | keyword_super | keyword_then | keyword_true | keyword_undef
                | keyword_when | keyword_yield | keyword_if | keyword_unless
                | keyword_while | keyword_until
                ;

arg		: lhs '=' lex_ctxt arg_rhs
                | var_lhs tOP_ASGN lex_ctxt arg_rhs
                | primary_value '[' opt_call_args rbracket tOP_ASGN lex_ctxt arg_rhs
                | primary_value call_op tIDENTIFIER tOP_ASGN lex_ctxt arg_rhs
                | primary_value call_op tCONSTANT tOP_ASGN lex_ctxt arg_rhs
                | primary_value tCOLON2 tIDENTIFIER tOP_ASGN lex_ctxt arg_rhs
                | primary_value tCOLON2 tCONSTANT tOP_ASGN lex_ctxt arg_rhs
                | tCOLON3 tCONSTANT tOP_ASGN lex_ctxt arg_rhs
                | backref tOP_ASGN lex_ctxt arg_rhs
                | arg tDOT2 arg
                | arg tDOT3 arg
                | arg tDOT2
                | arg tDOT3
                | tBDOT2 arg
                | tBDOT3 arg
                | arg '+' arg
                | arg '-' arg
                | arg '*' arg
                | arg '/' arg
                | arg '%' arg
                | arg tPOW arg
                | tUMINUS_NUM simple_numeric tPOW arg
                | tUPLUS arg
                | tUMINUS arg
                | arg '|' arg
                | arg '^' arg
                | arg '&' arg
                | arg tCMP arg
                | rel_expr   %prec tCMP
                | arg tEQ arg
                | arg tEQQ arg
                | arg tNEQ arg
                | arg tMATCH arg
                | arg tNMATCH arg
                | '!' arg
                | '~' arg
                | arg tLSHFT arg
                | arg tRSHFT arg
                | arg tANDOP arg
                | arg tOROP arg
                | keyword_defined '\n'? begin_defined arg
                | arg '?' arg '\n'? ':' arg
                | defn_head[head] f_opt_paren_args[args] '=' endless_arg[bodystmt]
                | defs_head[head] f_opt_paren_args[args] '=' endless_arg[bodystmt]
                | primary
                ;

endless_arg	: arg %prec modifier_rescue
                | endless_arg modifier_rescue after_rescue arg
                | keyword_not '\n'? endless_arg
                ;

relop		: '>'
                | '<'
                | tGEQ
                | tLEQ
                ;

rel_expr	: arg relop arg   %prec '>'
                | rel_expr relop arg   %prec '>'
                ;

lex_ctxt	: none
                ;

begin_defined	: lex_ctxt
                ;

after_rescue	: lex_ctxt
                ;

arg_value	: arg
                ;

aref_args	: none
                | args trailer
                | args ',' assocs trailer
                | assocs trailer
                ;

arg_rhs 	: arg   %prec tOP_ASGN
                | arg modifier_rescue after_rescue arg
                ;

paren_args	: '(' opt_call_args rparen
                | '(' args ',' args_forward rparen
                | '(' args_forward rparen
                ;

opt_paren_args	: none
                | paren_args
                ;

opt_call_args	: none
                | call_args
                | args ','
                | args ',' assocs ','
                | assocs ','
                ;

call_args	: command
                | args opt_block_arg
                | assocs opt_block_arg
                | args ',' assocs opt_block_arg
                | block_arg
                ;

command_args	: call_args
                ;

block_arg	: tAMPER arg_value
                | tAMPER
                ;

opt_block_arg	: ',' block_arg
                | none
                ;

args		: arg_value
                | arg_splat
                | args ',' arg_value
                | args ',' arg_splat
                ;

arg_splat	: tSTAR arg_value
                | tSTAR
                ;

mrhs_arg	: mrhs
                | arg_value
                ;

mrhs		: args ',' arg_value
                | args ',' tSTAR arg_value
                | tSTAR arg_value
                ;

primary		: literal
                | strings
                | xstring
                | regexp
                | words
                | qwords
                | symbols
                | qsymbols
                | var_ref
                | backref
                | tFID
                | k_begin
                  bodystmt
                  k_end
                | tLPAREN_ARG compstmt  ')'
                | tLPAREN compstmt ')'
                | primary_value tCOLON2 tCONSTANT
                | tCOLON3 tCONSTANT
                | tLBRACK aref_args ']'
                | tLBRACE assoc_list '}'
                | k_return
                | k_yield '(' call_args rparen
                | k_yield '(' rparen
                | k_yield
                | keyword_defined '\n'? '(' begin_defined expr rparen
                | keyword_not '(' expr rparen
                | keyword_not '(' rparen
                | fcall brace_block
                | method_call
                | method_call brace_block
                | lambda
                | k_if expr_value then
                  compstmt
                  if_tail
                  k_end
                | k_unless expr_value then
                  compstmt
                  opt_else
                  k_end
                | k_while expr_value_do
                  compstmt
                  k_end
                | k_until expr_value_do
                  compstmt
                  k_end
                | k_case expr_value terms?
                    <labels>
                  case_body
                  k_end
                | k_case terms?
                    <labels>
                  case_body
                  k_end
                | k_case expr_value terms?
                  p_case_body
                  k_end
                | k_for for_var keyword_in expr_value_do
                  compstmt
                  k_end
                | k_class cpath superclass
                  bodystmt
                  k_end
                | k_class tLSHFT expr_value
                  term
                  bodystmt
                  k_end
                | k_module cpath
                  bodystmt
                  k_end
                | defn_head[head]
                  f_arglist[args]
                  bodystmt
                  k_end
                | defs_head[head]
                  f_arglist[args]
                  bodystmt
                  k_end
                | keyword_break
                | keyword_next
                | keyword_redo
                | keyword_retry
                ;

primary_value	: primary
                ;

k_begin		: keyword_begin
                ;

k_if		: keyword_if
                ;

k_unless	: keyword_unless
                ;

k_while		: keyword_while allow_exits
                ;

k_until		: keyword_until allow_exits
                ;

k_case		: keyword_case
                ;

k_for		: keyword_for allow_exits
                ;

k_class		: keyword_class
                ;

k_module	: keyword_module
                ;

k_def		: keyword_def
                ;

k_do		: keyword_do
                ;

k_do_block	: keyword_do_block
                ;

k_rescue	: keyword_rescue
                ;

k_ensure	: keyword_ensure
                ;

k_when		: keyword_when
                ;

k_else		: keyword_else
                ;

k_elsif 	: keyword_elsif
                ;

k_end		: keyword_end
                | tDUMNY_END
                ;

k_return	: keyword_return
                ;

k_yield 	: keyword_yield
                ;

then		: term
                | keyword_then
                | term keyword_then
                ;

do		: term
                | keyword_do_cond
                ;

if_tail		: opt_else
                | k_elsif expr_value then
                  compstmt
                  if_tail
                ;

opt_else	: none
                | k_else compstmt
                ;

for_var		: lhs
                | mlhs
                ;

f_marg		: f_norm_arg
                | tLPAREN f_margs rparen
                ;

f_marg_list	: f_marg
                | f_marg_list ',' f_marg
                ;

f_margs		: f_marg_list
                | f_marg_list ',' f_rest_marg
                | f_marg_list ',' f_rest_marg ',' f_marg_list
                | f_rest_marg
                | f_rest_marg ',' f_marg_list
                ;

f_rest_marg	: tSTAR f_norm_arg
                | tSTAR
                ;

f_any_kwrest	: f_kwrest
                | f_no_kwarg
                ;

f_eq		:  '=';

block_args_tail	: f_kwarg(f_block_kw) ',' f_kwrest opt_f_block_arg
                | f_kwarg(f_block_kw) opt_f_block_arg
                | f_any_kwrest opt_f_block_arg
                | f_block_arg
                ;

excessed_comma	: ','
                ;

block_param	: f_arg ',' f_optarg(primary_value) ',' f_rest_arg opt_args_tail(block_args_tail)
                | f_arg ',' f_optarg(primary_value) ',' f_rest_arg ',' f_arg opt_args_tail(block_args_tail)
                | f_arg ',' f_optarg(primary_value) opt_args_tail(block_args_tail)
                | f_arg ',' f_optarg(primary_value) ',' f_arg opt_args_tail(block_args_tail)
                | f_arg ',' f_rest_arg opt_args_tail(block_args_tail)
                | f_arg excessed_comma
                | f_arg ',' f_rest_arg ',' f_arg opt_args_tail(block_args_tail)
                | f_arg opt_args_tail(block_args_tail)
                | f_optarg(primary_value) ',' f_rest_arg opt_args_tail(block_args_tail)
                | f_optarg(primary_value) ',' f_rest_arg ',' f_arg opt_args_tail(block_args_tail)
                | f_optarg(primary_value) opt_args_tail(block_args_tail)
                | f_optarg(primary_value) ',' f_arg opt_args_tail(block_args_tail)
                | f_rest_arg opt_args_tail(block_args_tail)
                | f_rest_arg ',' f_arg opt_args_tail(block_args_tail)
                | block_args_tail
                ;

opt_block_param	: none
                | block_param_def
                ;

block_param_def	: '|' opt_bv_decl '|'
                | '|' block_param opt_bv_decl '|'
                ;

opt_bv_decl	: '\n'?
                | '\n'? ';' bv_decls '\n'?
                ;

bv_decls	: bvar
                | bv_decls ',' bvar
                ;

bvar		: tIDENTIFIER
                | f_bad_arg
                ;

max_numparam	:
                ;

numparam	:
                ;

it_id           :
                ;

lambda		: tLAMBDA[lpar]
                    [dyna]<vars>
                  max_numparam numparam it_id allow_exits
                  f_larglist[args]
                  lambda_body[body]
                ;

f_larglist	: '(' f_args opt_bv_decl ')'
                | f_args
                ;

lambda_body	: tLAMBEG compstmt '}'
                | keyword_do_LAMBDA
                  bodystmt k_end
                ;

do_block	: k_do_block do_body k_end
                ;

block_call	: command do_block
                | block_call call_op2 operation2 opt_paren_args
                | block_call call_op2 operation2 opt_paren_args brace_block
                | block_call call_op2 operation2 command_args do_block
                ;

method_call	: fcall paren_args
                | primary_value call_op operation2 opt_paren_args
                | primary_value tCOLON2 operation2 paren_args
                | primary_value tCOLON2 operation3
                | primary_value call_op paren_args
                | primary_value tCOLON2 paren_args
                | keyword_super paren_args
                | keyword_super
                | primary_value '[' opt_call_args rbracket
                ;

brace_block	: '{' brace_body '}'
                | k_do do_body k_end
                ;

brace_body	: [dyna]<vars>
                  max_numparam numparam it_id allow_exits
                  opt_block_param[args] compstmt
                ;

do_body 	:   [dyna]<vars>
                  max_numparam numparam it_id allow_exits
                  opt_block_param[args] bodystmt
                ;

case_args	: arg_value
                | tSTAR arg_value
                | case_args ',' arg_value
                | case_args ',' tSTAR arg_value
                ;

case_body	: k_when case_args then
                  compstmt
                  cases
                ;

cases		: opt_else
                | case_body
                ;

p_pvtbl 	: ;
p_pktbl 	: ;

p_in_kwarg	:
                ;

p_case_body	: keyword_in
                  p_in_kwarg[ctxt] p_pvtbl p_pktbl
                  p_top_expr[expr] then
                  compstmt
                  p_cases[cases]
                ;

p_cases 	: opt_else
                | p_case_body
                ;

p_top_expr	: p_top_expr_body
                | p_top_expr_body modifier_if expr_value
                | p_top_expr_body modifier_unless expr_value
                ;

p_top_expr_body : p_expr
                | p_expr ','
                | p_expr ',' p_args
                | p_find
                | p_args_tail
                | p_kwargs
                ;

p_expr		: p_as
                ;

p_as		: p_expr tASSOC p_variable
                | p_alt
                ;

p_alt		: p_alt '|' p_expr_basic
                | p_expr_basic
                ;

p_lparen	: '(' p_pktbl
                ;

p_lbracket	: '[' p_pktbl
                ;

p_expr_basic	: p_value
                | p_variable
                | p_const p_lparen[p_pktbl] p_args rparen
                | p_const p_lparen[p_pktbl] p_find rparen
                | p_const p_lparen[p_pktbl] p_kwargs rparen
                | p_const '(' rparen
                | p_const p_lbracket[p_pktbl] p_args rbracket
                | p_const p_lbracket[p_pktbl] p_find rbracket
                | p_const p_lbracket[p_pktbl] p_kwargs rbracket
                | p_const '[' rbracket
                | tLBRACK p_args rbracket
                | tLBRACK p_find rbracket
                | tLBRACK rbracket
                | tLBRACE p_pktbl lex_ctxt[ctxt]
                  p_kwargs rbrace
                | tLBRACE rbrace
                | tLPAREN p_pktbl p_expr rparen
                ;

p_args		: p_expr
                | p_args_head
                | p_args_head p_arg
                | p_args_head p_rest
                | p_args_head p_rest ',' p_args_post
                | p_args_tail
                ;

p_args_head	: p_arg ','
                | p_args_head p_arg ','
                ;

p_args_tail	: p_rest
                | p_rest ',' p_args_post
                ;

p_find		: p_rest ',' p_args_post ',' p_rest
                ;

p_rest		: tSTAR tIDENTIFIER
                | tSTAR
                ;

p_args_post	: p_arg
                | p_args_post ',' p_arg
                ;

p_arg		: p_expr
                ;

p_kwargs	: p_kwarg ',' p_any_kwrest
                | p_kwarg
                | p_kwarg ','
                | p_any_kwrest
                ;

p_kwarg 	: p_kw
                | p_kwarg ',' p_kw
                ;

p_kw		: p_kw_label p_expr
                | p_kw_label
                ;

p_kw_label	: tLABEL
                | tSTRING_BEG string_contents tLABEL_END
                ;

p_kwrest	: kwrest_mark tIDENTIFIER
                | kwrest_mark
                ;

p_kwnorest	: kwrest_mark keyword_nil
                ;

p_any_kwrest	: p_kwrest
                | p_kwnorest
                ;

p_value 	: p_primitive
                | p_primitive tDOT2 p_primitive
                | p_primitive tDOT3 p_primitive
                | p_primitive tDOT2
                | p_primitive tDOT3
                | p_var_ref
                | p_expr_ref
                | p_const
                | tBDOT2 p_primitive
                | tBDOT3 p_primitive
                ;

p_primitive	: literal
                | strings
                | xstring
                | regexp
                | words
                | qwords
                | symbols
                | qsymbols
                | keyword_variable
                | lambda
                ;

p_variable	: tIDENTIFIER
                ;

p_var_ref	: '^' tIDENTIFIER
                | '^' nonlocal_var
                ;

p_expr_ref	: '^' tLPAREN expr_value rparen
                ;

p_const 	: tCOLON3 cname
                | p_const tCOLON2 cname
                | tCONSTANT
                ;

opt_rescue	: k_rescue exc_list exc_var then
                  compstmt
                  opt_rescue
                | none
                ;

exc_list	: arg_value
                | mrhs
                | none
                ;

exc_var		: tASSOC lhs
                | none
                ;

opt_ensure	: k_ensure compstmt
                | none
                ;

literal		: numeric
                | symbol
                ;

strings		: string
                ;

string		: tCHAR
                | string1
                | string string1
                ;

string1		: tSTRING_BEG string_contents tSTRING_END
                ;

xstring		: tXSTRING_BEG xstring_contents tSTRING_END
                ;

regexp		: tREGEXP_BEG regexp_contents tREGEXP_END
                ;

words		: words(tWORDS_BEG, word_list) <node>
                ;

word_list	:
                | word_list word ' '+
                ;

word		: string_content
                | word string_content
                ;

symbols 	: words(tSYMBOLS_BEG, symbol_list) <node>
                ;

symbol_list	:
                | symbol_list word ' '+
                ;

qwords		: words(tQWORDS_BEG, qword_list) <node>
                ;

qsymbols	: words(tQSYMBOLS_BEG, qsym_list) <node>
                ;

qword_list	:
                | qword_list tSTRING_CONTENT ' '+
                ;

qsym_list	:
                | qsym_list tSTRING_CONTENT ' '+
                ;

string_contents :
                | string_contents string_content
                ;

xstring_contents:
                | xstring_contents string_content
                ;

regexp_contents:
                | regexp_contents string_content
                ;

string_content	: tSTRING_CONTENT
                | tSTRING_DVAR
                    <strterm>
                  string_dvar
                | tSTRING_DBEG[state]
                    [term]<strterm>
                    [brace]<num>
                    [indent]<num>
                  compstmt string_dend
                ;

string_dend	: tSTRING_DEND
                | END_OF_INPUT
                ;

string_dvar	: nonlocal_var
                | backref
                ;

symbol		: ssym
                | dsym
                ;

ssym		: tSYMBEG sym
                ;

sym		: fname
                | nonlocal_var
                ;

dsym		: tSYMBEG string_contents tSTRING_END
                ;

numeric 	: simple_numeric
                | tUMINUS_NUM simple_numeric   %prec tLOWEST
                ;

simple_numeric	: tINTEGER
                | tFLOAT
                | tRATIONAL
                | tIMAGINARY
                ;

nonlocal_var    : tIVAR
                | tGVAR
                | tCVAR
                ;

user_variable	: ident_or_const
                | nonlocal_var
                ;

keyword_variable: keyword_nil
                | keyword_self
                | keyword_true
                | keyword_false
                | keyword__FILE__
                | keyword__LINE__
                | keyword__ENCODING__
                ;

var_ref		: user_variable
                | keyword_variable
                ;

var_lhs		: user_variable
                | keyword_variable
                ;

backref		: tNTH_REF
                | tBACK_REF
                ;

superclass	: '<'
                  expr_value term
                |
                ;

f_opt_paren_args: f_paren_args
                | none
                ;

f_paren_args	: '(' f_args rparen
                ;

f_arglist	: f_paren_args
                |   <ctxt>
                  f_args term
                ;

args_tail	: f_kwarg(f_kw) ',' f_kwrest opt_f_block_arg
                | f_kwarg(f_kw) opt_f_block_arg
                | f_any_kwrest opt_f_block_arg
                | f_block_arg
                | args_forward
                ;

f_args		: f_arg ',' f_optarg(arg_value) ',' f_rest_arg opt_args_tail(args_tail)
                | f_arg ',' f_optarg(arg_value) ',' f_rest_arg ',' f_arg opt_args_tail(args_tail)
                | f_arg ',' f_optarg(arg_value) opt_args_tail(args_tail)
                | f_arg ',' f_optarg(arg_value) ',' f_arg opt_args_tail(args_tail)
                | f_arg ',' f_rest_arg opt_args_tail(args_tail)
                | f_arg ',' f_rest_arg ',' f_arg opt_args_tail(args_tail)
                | f_arg opt_args_tail(args_tail)
                | f_optarg(arg_value) ',' f_rest_arg opt_args_tail(args_tail)
                | f_optarg(arg_value) ',' f_rest_arg ',' f_arg opt_args_tail(args_tail)
                | f_optarg(arg_value) opt_args_tail(args_tail)
                | f_optarg(arg_value) ',' f_arg opt_args_tail(args_tail)
                | f_rest_arg opt_args_tail(args_tail)
                | f_rest_arg ',' f_arg opt_args_tail(args_tail)
                | args_tail
                |
                ;

args_forward	: tBDOT3
                ;

f_bad_arg	: tCONSTANT
                | tIVAR
                | tGVAR
                | tCVAR
                ;

f_norm_arg	: f_bad_arg
                | tIDENTIFIER
                ;

f_arg_asgn	: f_norm_arg
                ;

f_arg_item	: f_arg_asgn
                | tLPAREN f_margs rparen
                ;

f_arg		: f_arg_item
                | f_arg ',' f_arg_item
                ;

f_label 	: tLABEL
                ;

f_kw		: f_label arg_value
                | f_label
                ;

f_block_kw	: f_label primary_value
                | f_label
                ;

kwrest_mark	: tPOW
                | tDSTAR
                ;

f_no_kwarg	: p_kwnorest
                ;

f_kwrest	: kwrest_mark tIDENTIFIER
                | kwrest_mark
                ;

restarg_mark	: '*'
                | tSTAR
                ;

f_rest_arg	: restarg_mark tIDENTIFIER
                | restarg_mark
                ;

blkarg_mark	: '&'
                | tAMPER
                ;

f_block_arg	: blkarg_mark tIDENTIFIER
                | blkarg_mark
                ;

opt_f_block_arg	: ',' f_block_arg
                | none
                ;

singleton	: var_ref
                | '('  expr rparen
                ;

assoc_list	: none
                | assocs trailer
                ;

assocs		: assoc
                | assocs ',' assoc
                ;

assoc		: arg_value tASSOC arg_value
                | tLABEL arg_value
                | tLABEL
                | tSTRING_BEG string_contents tLABEL_END arg_value
                | tDSTAR arg_value
                | tDSTAR
                ;

operation	: ident_or_const
                | tFID
                ;

operation2	: operation
                | op
                ;

operation3	: tIDENTIFIER
                | tFID
                | op
                ;

dot_or_colon	: '.'
                | tCOLON2
                ;

call_op 	: '.'
                | tANDDOT
                ;

call_op2	: call_op
                | tCOLON2
                ;

rparen		: '\n'? ')'
                ;

rbracket	: '\n'? ']'
                ;

rbrace		: '\n'? '}'
                ;

trailer		: '\n'?
                | ','
                ;

term		: ';'
                | '\n'
                ;

terms		: term
                | terms ';'
                ;

none		:
                ;
