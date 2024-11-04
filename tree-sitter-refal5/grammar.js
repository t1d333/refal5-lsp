/**
 * @file Refal5 grammar for tree-sitter
 * @author Kirill Kiselev <kirillkiselev2003@yandex.ru>
 * @license MIT
 */

/// <reference types="tree-sitter-cli/dsl" />
// @ts-check
//
// Grammar from http://refal.botik.ru/book/html/
// expression ::= empty 
//                term expression 
//
// term ::= symbol 
//          variable  
//          (expression)
//          
//
// f-name ::= identifier 
// empty ::= 
// A pattern expression, or pattern, is an expression which includes no activation brackets. A ground expression is an expression which includes no variables.
//
// 3. Sentences and programs
// A Refal program is a list of function definitions (f-definition's — see below) and external function declarations (external-decl). Semicolons must be used to separate an external declaration from the following function definition; they may also separate function definitions.
//
// program ::= f-definition 
//             f-definition  program 
//             f-definition ;  program 
//             external-decl ;  program 
//             program  external-decl ;
//
// f-definition ::= f-name {  block }
//              $ENTRY  f-name {  block }
//
// external-decl ::= $EXTERNAL  f-name-list 
//                   $EXTERN  f-name-list 
//                   $EXTRN  f-name-list 
//
// f-name-list ::= f-name 
//                 f-name , f-name-list 
//
// f-name ::= identifier 
//
// block ::= sentence 
//           sentence ;
//           sentence ;  block 
//
// sentence ::= left-side  conditions  =   right-side 
//              left-side  conditions  ,  block-ending 
//
// left-side ::= pattern 
//
// conditions ::= empty 
//                , arg : pattern conditions 
//
// arg ::= expression 
//
// right-side ::= expression 
//
// block-ending ::= arg : { block }

module.exports = grammar({
  name: "refal5",

  rules: {
    source_file: $ => repeat(
      choice(
        $.external_declaration,
        $.function_definition,
      )
    ),
    
    external_declaration: $ => seq(
      $.external_modifier,
      field('func_name_list', $.function_name_list),
      ';'
    ),

    function_name_list: $ => repeat1(
      $.ident
    ),
    
    function_definition: $ => seq(
      optional($.entry_modifier),
      field('func_name', $.ident,),
      $.body,
    ),

    body: $ => seq(
      '{',
      repeat(
        $.sentence
      ),
      '}'
    ),

    sentence: $ => choice(
      $.sentenceEq, 
      $.sentenceBlock 
    ),

    sentenceEq: $ => seq(
      optional($.pattern), 
      repeat($.condition),
      "=",
      optional($.expr),
      ';'
    ),

    sentenceBlock: $ => seq(
      optional($.pattern), 
      repeat($.condition),
      ',',
      $.sentenceBlockEnd,
    ),
    
    sentenceBlockEnd: $ => seq(
      optional($.expr), 
      ':',
      '{',
      $.body, 
      '}'
    ),
    
    pattern: $ => repeat1(choice(
      $.ident,
      $.string,
      $.number,
      $.variable,
      seq('(', optional($.pattern), ')')
    )),

    condition: $ => seq(
        ',',
        optional($.expr),
        ':',
        optional($.pattern),
    ),

    expr: $ => repeat1(
      choice(
        $.ident,
        $.string,
        $.number,
        $.variable,
        seq('<', $.ident,  optional($.expr), '>'),
        seq('(', optional($.expr), ')')
      )
    ),

    variable: $ => seq(
      field('type', $.type),
      '.',
      field('name', $.ident),
    ), 
    
    ident : $ => /(([A-Za-z][A-Za-z0-9_-]*)|([0-9]+))/, 

    type: $ => choice(
      's',
      'e',
      't'
    ),

    string: $ => /'[^\n']*'|\"[^\n\"]*\"/,
    number: $ => /\d+/,
    entry_modifier: $ => /\$ENTRY/,
    external_modifier: $ => choice(
        '$EXTERNAL',
        '$EXTERN',
        '$EXTRN'
    )
  }
});
