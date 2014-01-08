//line cc.y:34
package cc

import __yyfmt__ "fmt"

//line cc.y:34
type typeClass struct {
	c TypeStorage
	q TypeQual
	t *Type
}

type idecor struct {
	d func(*Type) (*Type, string)
	i *Init
}

//line cc.y:49
type yySymType struct {
	yys      int
	abdecor  func(*Type) *Type
	decl     *Decl
	decls    []*Decl
	decor    func(*Type) (*Type, string)
	decors   []func(*Type) (*Type, string)
	expr     *Expr
	exprs    []*Expr
	idec     idecor
	idecs    []idecor
	init     *Init
	inits    []*Init
	label    *Label
	labels   []*Label
	span     Span
	prefix   *Prefix
	prefixes []*Prefix
	stmt     *Stmt
	stmts    []*Stmt
	str      string
	strs     []string
	tc       typeClass
	tk       TypeKind
	typ      *Type
}

const tokARGBEGIN = 57346
const tokARGEND = 57347
const tokAuto = 57348
const tokBreak = 57349
const tokCase = 57350
const tokChar = 57351
const tokConst = 57352
const tokContinue = 57353
const tokDefault = 57354
const tokDo = 57355
const tokDotDotDot = 57356
const tokDouble = 57357
const tokEnum = 57358
const tokError = 57359
const tokExtern = 57360
const tokFloat = 57361
const tokFor = 57362
const tokGoto = 57363
const tokIf = 57364
const tokInline = 57365
const tokInt = 57366
const tokLitChar = 57367
const tokLong = 57368
const tokName = 57369
const tokNumber = 57370
const tokOffsetof = 57371
const tokRegister = 57372
const tokReturn = 57373
const tokShort = 57374
const tokSigned = 57375
const tokStatic = 57376
const tokStruct = 57377
const tokSwitch = 57378
const tokTypeName = 57379
const tokTypedef = 57380
const tokUnion = 57381
const tokUnsigned = 57382
const tokVaArg = 57383
const tokVoid = 57384
const tokVolatile = 57385
const tokWhile = 57386
const tokString = 57387
const tokShift = 57388
const tokElse = 57389
const tokAddEq = 57390
const tokSubEq = 57391
const tokMulEq = 57392
const tokDivEq = 57393
const tokModEq = 57394
const tokLshEq = 57395
const tokRshEq = 57396
const tokAndEq = 57397
const tokXorEq = 57398
const tokOrEq = 57399
const tokOrOr = 57400
const tokAndAnd = 57401
const tokEqEq = 57402
const tokNotEq = 57403
const tokLtEq = 57404
const tokGtEq = 57405
const tokLsh = 57406
const tokRsh = 57407
const tokCast = 57408
const tokSizeof = 57409
const tokUnary = 57410
const tokDec = 57411
const tokInc = 57412
const tokArrow = 57413
const startExpr = 57414
const startProg = 57415
const tokEOF = 57416

var yyToknames = []string{
	"tokARGBEGIN",
	"tokARGEND",
	"tokAuto",
	"tokBreak",
	"tokCase",
	"tokChar",
	"tokConst",
	"tokContinue",
	"tokDefault",
	"tokDo",
	"tokDotDotDot",
	"tokDouble",
	"tokEnum",
	"tokError",
	"tokExtern",
	"tokFloat",
	"tokFor",
	"tokGoto",
	"tokIf",
	"tokInline",
	"tokInt",
	"tokLitChar",
	"tokLong",
	"tokName",
	"tokNumber",
	"tokOffsetof",
	"tokRegister",
	"tokReturn",
	"tokShort",
	"tokSigned",
	"tokStatic",
	"tokStruct",
	"tokSwitch",
	"tokTypeName",
	"tokTypedef",
	"tokUnion",
	"tokUnsigned",
	"tokVaArg",
	"tokVoid",
	"tokVolatile",
	"tokWhile",
	"tokString",
	"tokShift",
	"tokElse",
	" {",
	" ,",
	" =",
	"tokAddEq",
	"tokSubEq",
	"tokMulEq",
	"tokDivEq",
	"tokModEq",
	"tokLshEq",
	"tokRshEq",
	"tokAndEq",
	"tokXorEq",
	"tokOrEq",
	" ?",
	" :",
	"tokOrOr",
	"tokAndAnd",
	" |",
	" ^",
	" &",
	"tokEqEq",
	"tokNotEq",
	" <",
	" >",
	"tokLtEq",
	"tokGtEq",
	"tokLsh",
	"tokRsh",
	" +",
	" -",
	" *",
	" /",
	" %",
	"tokCast",
	" !",
	" ~",
	"tokSizeof",
	"tokUnary",
	" .",
	" [",
	" ]",
	" (",
	" )",
	"tokDec",
	"tokInc",
	"tokArrow",
	"startExpr",
	"startProg",
	"tokEOF",
}
var yyStatenames = []string{}

const yyEofCode = 1
const yyErrCode = 2
const yyMaxDepth = 200

//line yacctab:1
var yyExca = []int{
	-1, 1,
	1, -1,
	-2, 0,
	-1, 114,
	49, 97,
	98, 97,
	-2, 174,
	-1, 130,
	48, 165,
	-2, 139,
	-1, 134,
	48, 165,
	-2, 144,
	-1, 231,
	98, 200,
	-2, 164,
	-1, 259,
	62, 132,
	-2, 88,
}

const yyNprod = 210
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 1415

var yyAct = []int{

	7, 293, 111, 260, 220, 327, 278, 29, 217, 256,
	270, 190, 110, 97, 98, 99, 100, 101, 102, 103,
	104, 105, 233, 207, 328, 47, 211, 230, 5, 108,
	187, 215, 292, 209, 219, 119, 127, 27, 125, 130,
	134, 114, 356, 123, 354, 344, 339, 109, 337, 322,
	321, 319, 286, 279, 181, 298, 282, 237, 58, 30,
	136, 137, 138, 139, 140, 141, 142, 143, 144, 145,
	146, 147, 148, 149, 150, 151, 152, 153, 154, 116,
	156, 157, 158, 159, 160, 161, 162, 163, 164, 165,
	166, 124, 6, 32, 3, 2, 184, 359, 170, 171,
	285, 95, 91, 155, 90, 343, 93, 92, 94, 353,
	347, 259, 183, 180, 182, 228, 169, 204, 346, 33,
	242, 133, 345, 275, 121, 116, 129, 183, 274, 182,
	248, 109, 244, 183, 176, 182, 172, 173, 177, 199,
	197, 249, 330, 189, 175, 174, 329, 326, 324, 179,
	122, 206, 128, 65, 66, 60, 61, 62, 63, 64,
	192, 112, 117, 193, 191, 95, 91, 107, 90, 273,
	93, 92, 94, 118, 289, 201, 117, 250, 198, 204,
	124, 333, 332, 168, 216, 218, 265, 118, 223, 281,
	196, 280, 262, 245, 200, 213, 59, 235, 225, 226,
	201, 236, 186, 189, 195, 216, 194, 231, 205, 202,
	271, 272, 355, 96, 227, 129, 335, 56, 221, 213,
	129, 132, 224, 258, 234, 247, 287, 261, 251, 240,
	31, 133, 1, 35, 202, 225, 239, 246, 243, 218,
	241, 128, 231, 257, 122, 11, 128, 268, 188, 126,
	57, 46, 60, 61, 62, 63, 64, 253, 213, 131,
	135, 120, 95, 91, 301, 90, 284, 93, 92, 94,
	269, 295, 276, 291, 189, 290, 302, 238, 267, 277,
	288, 223, 297, 113, 115, 283, 167, 263, 218, 226,
	296, 264, 257, 254, 255, 232, 299, 229, 26, 4,
	304, 240, 210, 185, 28, 178, 0, 0, 0, 0,
	0, 323, 0, 320, 0, 325, 0, 331, 0, 0,
	305, 0, 0, 223, 62, 63, 64, 0, 0, 0,
	336, 0, 95, 91, 0, 90, 0, 93, 92, 94,
	0, 0, 0, 0, 0, 0, 0, 350, 351, 352,
	349, 338, 0, 0, 340, 341, 0, 358, 0, 0,
	357, 360, 0, 0, 0, 0, 0, 0, 348, 79,
	80, 81, 82, 83, 84, 85, 86, 87, 88, 89,
	78, 342, 77, 76, 75, 74, 73, 71, 72, 67,
	68, 69, 70, 65, 66, 60, 61, 62, 63, 64,
	0, 0, 0, 0, 0, 95, 91, 0, 90, 0,
	93, 92, 94, 79, 80, 81, 82, 83, 84, 85,
	86, 87, 88, 89, 78, 0, 77, 76, 75, 74,
	73, 71, 72, 67, 68, 69, 70, 65, 66, 60,
	61, 62, 63, 64, 0, 0, 0, 0, 0, 95,
	91, 300, 90, 0, 93, 92, 94, 79, 80, 81,
	82, 83, 84, 85, 86, 87, 88, 89, 78, 0,
	77, 76, 75, 74, 73, 71, 72, 67, 68, 69,
	70, 65, 66, 60, 61, 62, 63, 64, 0, 0,
	0, 0, 0, 95, 91, 0, 90, 266, 93, 92,
	94, 208, 79, 80, 81, 82, 83, 84, 85, 86,
	87, 88, 89, 78, 0, 77, 76, 75, 74, 73,
	71, 72, 67, 68, 69, 70, 65, 66, 60, 61,
	62, 63, 64, 0, 0, 0, 0, 0, 95, 91,
	0, 90, 0, 93, 92, 94, 79, 80, 81, 82,
	83, 84, 85, 86, 87, 88, 89, 78, 0, 77,
	76, 75, 74, 73, 71, 72, 67, 68, 69, 70,
	65, 66, 60, 61, 62, 63, 64, 0, 0, 0,
	0, 0, 95, 91, 0, 90, 0, 93, 92, 94,
	50, 0, 0, 37, 56, 0, 0, 0, 0, 44,
	36, 0, 52, 43, 0, 0, 0, 55, 39, 10,
	40, 8, 9, 21, 54, 0, 38, 41, 51, 48,
	0, 34, 53, 49, 42, 23, 45, 57, 0, 24,
	76, 75, 74, 73, 71, 72, 67, 68, 69, 70,
	65, 66, 60, 61, 62, 63, 64, 0, 0, 0,
	0, 13, 95, 91, 0, 90, 0, 93, 92, 94,
	14, 15, 12, 0, 0, 0, 16, 17, 20, 0,
	0, 0, 0, 22, 306, 19, 18, 307, 316, 0,
	0, 308, 317, 309, 0, 0, 0, 0, 0, 0,
	310, 311, 312, 0, 0, 10, 0, 318, 9, 21,
	0, 313, 0, 0, 0, 0, 314, 0, 0, 0,
	0, 23, 0, 0, 315, 24, 0, 0, 222, 73,
	71, 72, 67, 68, 69, 70, 65, 66, 60, 61,
	62, 63, 64, 0, 0, 0, 0, 13, 95, 91,
	0, 90, 0, 93, 92, 94, 14, 15, 12, 0,
	0, 0, 16, 17, 20, 0, 0, 0, 0, 22,
	50, 19, 18, 37, 56, 0, 0, 0, 303, 44,
	36, 0, 52, 43, 0, 0, 0, 55, 39, 0,
	40, 0, 0, 0, 54, 0, 38, 41, 51, 48,
	0, 34, 53, 49, 42, 50, 45, 57, 37, 56,
	0, 0, 0, 0, 44, 36, 0, 52, 43, 0,
	0, 0, 55, 39, 0, 40, 0, 0, 0, 54,
	0, 38, 41, 51, 48, 0, 34, 53, 49, 42,
	0, 45, 57, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 294, 78, 0, 77, 76, 75, 74, 73, 71,
	72, 67, 68, 69, 70, 65, 66, 60, 61, 62,
	63, 64, 0, 0, 0, 0, 0, 95, 91, 0,
	90, 0, 93, 92, 94, 50, 252, 0, 37, 56,
	0, 0, 0, 0, 44, 36, 0, 52, 43, 0,
	0, 0, 55, 39, 0, 40, 0, 0, 0, 54,
	0, 38, 41, 51, 48, 0, 34, 53, 49, 42,
	0, 45, 57, 75, 74, 73, 71, 72, 67, 68,
	69, 70, 65, 66, 60, 61, 62, 63, 64, 0,
	0, 0, 0, 0, 95, 91, 0, 90, 0, 93,
	92, 94, 0, 74, 73, 71, 72, 67, 68, 69,
	70, 65, 66, 60, 61, 62, 63, 64, 0, 0,
	0, 0, 0, 95, 91, 25, 90, 0, 93, 92,
	94, 71, 72, 67, 68, 69, 70, 65, 66, 60,
	61, 62, 63, 64, 10, 0, 8, 9, 21, 95,
	91, 0, 90, 0, 93, 92, 94, 0, 0, 0,
	23, 0, 0, 0, 24, 0, 0, 203, 0, 0,
	72, 67, 68, 69, 70, 65, 66, 60, 61, 62,
	63, 64, 0, 0, 0, 0, 13, 95, 91, 0,
	90, 0, 93, 92, 94, 14, 15, 12, 0, 0,
	0, 16, 17, 20, 0, 271, 272, 0, 22, 0,
	19, 18, 67, 68, 69, 70, 65, 66, 60, 61,
	62, 63, 64, 10, 0, 8, 9, 21, 95, 91,
	0, 90, 50, 93, 92, 94, 56, 0, 0, 23,
	0, 0, 0, 24, 52, 0, 203, 0, 0, 55,
	10, 0, 8, 9, 21, 0, 54, 0, 0, 0,
	51, 0, 0, 0, 53, 13, 23, 0, 0, 57,
	24, 0, 0, 0, 14, 15, 12, 0, 0, 0,
	16, 17, 20, 0, 0, 0, 0, 22, 0, 19,
	18, 0, 13, 0, 0, 10, 0, 8, 9, 21,
	0, 14, 15, 12, 0, 0, 0, 16, 17, 20,
	0, 23, 0, 0, 22, 24, 19, 18, 10, 0,
	8, 9, 21, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 23, 0, 0, 13, 24, 0,
	0, 203, 0, 0, 0, 0, 14, 15, 12, 0,
	0, 0, 16, 17, 20, 0, 0, 0, 0, 106,
	0, 19, 18, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 16, 17, 20, 0, 0,
	0, 0, 22, 50, 19, 18, 37, 56, 0, 0,
	0, 0, 44, 36, 0, 52, 43, 0, 0, 0,
	55, 39, 0, 40, 0, 0, 0, 54, 0, 38,
	41, 51, 48, 0, 34, 53, 49, 42, 50, 45,
	57, 37, 56, 0, 0, 222, 214, 44, 36, 0,
	52, 43, 0, 0, 0, 55, 39, 0, 40, 212,
	0, 0, 54, 0, 38, 41, 51, 48, 0, 34,
	53, 49, 42, 0, 45, 57, 334, 50, 0, 0,
	37, 56, 0, 0, 0, 0, 44, 36, 0, 52,
	43, 0, 0, 0, 55, 39, 0, 40, 0, 0,
	0, 54, 0, 38, 41, 51, 48, 0, 34, 53,
	49, 42, 50, 45, 57, 37, 56, 0, 0, 0,
	0, 44, 36, 0, 52, 43, 0, 0, 0, 55,
	39, 0, 40, 0, 0, 0, 54, 0, 38, 41,
	51, 48, 0, 34, 53, 49, 42, 50, 45, 57,
	37, 56, 0, 0, 0, 0, 44, 0, 0, 52,
	43, 0, 0, 0, 55, 39, 0, 40, 0, 0,
	0, 54, 0, 38, 41, 51, 0, 0, 0, 53,
	0, 42, 0, 45, 57,
}
var yyPact = []int{

	0, -1000, -1000, 1075, 879, -38, 147, 496, -1000, -1000,
	-1000, 168, 1075, 1075, 1075, 1075, 1075, 1075, 1075, 1075,
	1120, 78, 584, 72, -1000, -1000, -1000, -1000, -1000, 98,
	1336, 1076, 1371, -1000, -1000, 194, 194, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, 1075,
	1075, 1075, 1075, 1075, 1075, 1075, 1075, 1075, 1075, 1075,
	1075, 1075, 1075, 1075, 1075, 1075, 1075, 1075, 1075, 1075,
	1075, 1075, 1075, 1075, 1075, 1075, 1075, 1075, 1075, 1075,
	1075, 1075, -1000, -1000, 194, 194, -1000, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 584, 1336, 55, 54,
	60, -1000, 1075, -44, 46, 153, -1000, 207, 98, -1000,
	1076, 1371, -1000, -1000, 1076, -1000, 1371, -1000, -1000, -1000,
	-1000, 158, -1000, -1000, -1000, 156, 496, 246, 246, 15,
	15, 15, 176, 176, 79, 79, 79, 79, 951, 992,
	913, 652, 887, 858, 566, 128, 496, 496, 496, 496,
	496, 496, 496, 496, 496, 496, 496, 50, 147, 90,
	-1000, -1000, 49, 145, 1048, -1000, 92, 207, 62, 60,
	452, -1000, 1262, 1075, 1048, 1227, 98, 98, 207, -1000,
	25, -1000, -1000, -1000, 1336, 197, 1075, -1000, -1000, 1143,
	1075, 15, -1000, -40, 1075, 60, 1262, 30, 1336, 42,
	144, -1000, -1000, 52, -1000, 89, 496, -1000, 496, -1000,
	-1000, -1000, -1000, 98, -1000, 46, 40, -1000, -1000, 789,
	-1000, 84, 143, -1000, 136, 791, 407, -1000, 969, 81,
	92, 38, -1000, 33, -1000, 1262, 92, 40, 207, 52,
	-1000, -1000, -1000, -1000, -45, 142, -1000, 40, 127, -1000,
	-1000, -41, 197, -1000, -1000, 1075, -1000, 3, -1000, 124,
	-1000, 194, 1075, -1000, -1000, -1000, -1000, 52, 754, -1000,
	84, 1075, -1000, -1000, 496, -1000, -42, 1048, -1000, -1000,
	-1000, 363, -1000, -1000, -1000, 670, -1000, 496, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -47, -1000, -48, -49, -1000,
	59, 194, 58, 1075, 57, 53, 1075, 120, 119, -1000,
	1301, -1000, -1000, 172, 1075, -50, 1075, -52, -1000, 1075,
	1075, 319, -1000, -1000, -1000, 16, -53, -1000, 32, -1000,
	28, 20, -1000, 1075, 1075, -1000, -1000, -1000, 19, -54,
	165, -1000, -1000, -56, 1075, -1000, -1000, 7, -1000, -1000,
	-1000,
}
var yyPgo = []int{

	0, 23, 305, 26, 304, 22, 32, 303, 302, 33,
	299, 298, 27, 297, 295, 11, 9, 294, 293, 0,
	31, 24, 5, 291, 287, 92, 286, 35, 284, 283,
	8, 278, 34, 277, 276, 271, 10, 270, 264, 4,
	1, 6, 251, 25, 93, 119, 36, 3, 223, 59,
	43, 249, 38, 248, 30, 245, 2, 233, 29, 12,
	230, 232, 228, 227, 226,
}
var yyR1 = []int{

	0, 61, 61, 10, 10, 21, 19, 19, 19, 19,
	19, 19, 19, 19, 19, 19, 19, 19, 19, 19,
	19, 19, 19, 19, 19, 19, 19, 19, 19, 19,
	19, 19, 19, 19, 19, 19, 19, 19, 19, 19,
	19, 19, 19, 19, 19, 19, 19, 19, 19, 19,
	19, 19, 19, 19, 19, 19, 19, 19, 19, 41,
	41, 41, 62, 39, 34, 34, 34, 40, 38, 38,
	38, 38, 38, 38, 38, 38, 38, 38, 38, 38,
	38, 38, 1, 1, 1, 2, 2, 2, 15, 15,
	15, 15, 15, 3, 3, 3, 3, 27, 27, 42,
	42, 42, 42, 42, 42, 43, 43, 44, 44, 44,
	44, 44, 44, 44, 44, 44, 45, 45, 46, 46,
	60, 56, 56, 56, 56, 56, 59, 58, 6, 11,
	11, 4, 47, 47, 57, 57, 16, 16, 12, 60,
	60, 36, 19, 19, 60, 60, 5, 23, 30, 30,
	32, 32, 32, 33, 33, 31, 31, 36, 64, 64,
	63, 63, 37, 37, 48, 48, 22, 22, 20, 20,
	25, 25, 26, 26, 7, 7, 35, 35, 8, 8,
	9, 9, 28, 28, 29, 29, 53, 53, 54, 54,
	49, 49, 50, 50, 51, 51, 52, 52, 17, 17,
	18, 18, 13, 13, 24, 24, 14, 14, 55, 55,
}
var yyR2 = []int{

	0, 3, 3, 0, 2, 1, 1, 1, 1, 1,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 5, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	2, 2, 2, 2, 2, 2, 2, 2, 2, 4,
	6, 4, 4, 3, 4, 4, 2, 2, 6, 0,
	2, 2, 0, 4, 3, 2, 2, 2, 1, 1,
	2, 3, 2, 2, 7, 9, 3, 5, 7, 3,
	5, 5, 0, 3, 1, 4, 4, 3, 1, 3,
	3, 4, 4, 1, 2, 2, 1, 1, 3, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 3, 3, 2, 2, 1, 2, 3, 1,
	1, 4, 1, 1, 1, 1, 1, 3, 3, 2,
	5, 2, 3, 3, 2, 6, 2, 2, 1, 1,
	2, 4, 5, 0, 3, 1, 3, 3, 0, 1,
	0, 1, 1, 2, 0, 1, 0, 1, 0, 1,
	1, 3, 0, 1, 0, 2, 0, 2, 1, 3,
	0, 1, 1, 3, 0, 1, 1, 2, 0, 1,
	1, 2, 0, 1, 1, 2, 0, 1, 1, 3,
	0, 1, 1, 2, 0, 1, 1, 3, 1, 2,
}
var yyChk = []int{

	-1000, -61, 95, 94, -10, -21, -25, -19, 27, 28,
	25, -55, 78, 67, 76, 77, 82, 83, 92, 91,
	84, 29, 89, 41, 45, 96, -11, -6, -4, -56,
	-49, -60, -44, -45, 37, -57, 16, 9, 32, 24,
	26, 33, 40, 19, 15, 42, -42, -43, 35, 39,
	6, 34, 18, 38, 30, 23, 10, 43, 96, 49,
	76, 77, 78, 79, 80, 74, 75, 70, 71, 72,
	73, 68, 69, 67, 66, 65, 64, 63, 61, 50,
	51, 52, 53, 54, 55, 56, 57, 58, 59, 60,
	89, 87, 92, 91, 93, 86, 45, -19, -19, -19,
	-19, -19, -19, -19, -19, -19, 89, 89, -58, -21,
	-59, -56, 89, -29, -15, -28, 27, 78, 89, -27,
	-60, -44, -45, -50, -49, -52, -51, -46, -45, -44,
	-47, -48, 27, 37, -47, -48, -19, -19, -19, -19,
	-19, -19, -19, -19, -19, -19, -19, -19, -19, -19,
	-19, -19, -19, -19, -19, -21, -19, -19, -19, -19,
	-19, -19, -19, -19, -19, -19, -19, -26, -25, -21,
	-47, -47, -58, -58, 90, 90, -1, 78, -2, 89,
	-19, 98, 89, 87, 50, -7, 49, -54, -53, -43,
	-15, -50, -52, -46, 48, 48, 62, 90, 88, 90,
	49, -19, -32, 48, 87, -54, 89, -1, 49, -9,
	-8, -3, 27, -59, 14, -20, -19, -30, -19, -32,
	-39, -6, 48, -56, -27, -15, -15, -43, 90, -13,
	-12, -59, -14, -5, 27, -19, -19, 97, -33, -20,
	-1, -9, 90, -58, 90, 49, -1, -15, 78, 89,
	88, -62, 97, -12, -18, -17, -16, -15, -48, 27,
	-47, -63, 49, -24, -23, 50, 90, -31, -30, -37,
	-36, 86, 87, 88, 90, 90, -3, -54, -41, 98,
	49, 62, 97, -5, -19, 97, 49, -64, -36, 50,
	-47, -19, -6, -40, 97, -35, -16, -19, 97, -30,
	88, -38, -34, 98, -39, -21, 4, 7, 11, 13,
	20, 21, 22, 31, 36, 44, 8, 12, 27, 98,
	-41, 98, 98, -40, 89, -47, 89, -22, -21, 89,
	89, -19, 62, 62, 5, 44, -22, 98, -21, 98,
	-21, -21, 62, 89, 98, 90, 90, 90, -21, -22,
	-40, -40, -40, 90, 98, 47, 98, -22, -40, 90,
	-40,
}
var yyDef = []int{

	0, -2, 3, 0, 0, 0, 5, 170, 6, 7,
	8, 9, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 208, 1, 4, 129, 130, 184,
	121, 192, 196, 190, 120, 164, 164, 107, 108, 109,
	110, 111, 112, 113, 114, 115, 116, 117, 134, 135,
	99, 100, 101, 102, 103, 104, 105, 106, 2, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	172, 0, 56, 57, 0, 0, 209, 40, 41, 42,
	43, 44, 45, 46, 47, 48, 0, 0, 0, 0,
	82, 126, 0, 0, -2, 185, 88, 188, 0, 182,
	192, 196, 191, 124, 193, 125, 197, 194, 118, 119,
	-2, 0, 132, 133, -2, 0, 171, 10, 11, 12,
	13, 14, 15, 16, 17, 18, 19, 20, 21, 22,
	23, 24, 25, 26, 27, 0, 29, 30, 31, 32,
	33, 34, 35, 36, 37, 38, 39, 0, 173, 0,
	142, 143, 0, 0, 0, 53, 127, 188, 84, 82,
	0, 128, 180, 168, 0, 0, 0, 0, 189, 186,
	0, 122, 123, 195, 0, 0, 0, 54, 55, 49,
	0, 51, 52, 153, 168, 82, 180, 0, 0, 0,
	181, 178, 93, 82, 96, 0, 169, 98, 148, 149,
	131, 175, 62, 184, 183, 97, 89, 187, 90, 0,
	202, -2, 160, 206, 204, 28, 0, 150, 0, 0,
	83, 0, 87, 0, 91, 0, 94, 95, 188, 82,
	92, 59, 140, 203, 0, 201, 198, 136, 0, -2,
	165, 0, 161, 146, 205, 0, 50, 0, 155, 158,
	162, 0, 0, 86, 85, 58, 179, 82, 176, 138,
	164, 0, 145, 207, 147, 151, 154, 0, 163, 159,
	141, 0, 60, 61, 63, 0, 199, 137, 152, 156,
	157, 67, 177, 68, 69, 0, 59, 0, 0, 176,
	0, 0, 0, 166, 0, 0, 0, 0, 6, 70,
	176, 72, 73, 0, 166, 0, 0, 0, 167, 0,
	0, 0, 65, 66, 71, 0, 0, 76, 0, 79,
	0, 0, 64, 0, 166, 176, 176, 176, 0, 0,
	77, 80, 81, 0, 166, 176, 74, 0, 78, 176,
	75,
}
var yyTok1 = []int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 82, 3, 3, 3, 80, 67, 3,
	89, 90, 78, 76, 49, 77, 86, 79, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 62, 98,
	70, 50, 71, 61, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 87, 3, 88, 66, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 48, 65, 97, 83,
}
var yyTok2 = []int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 36, 37, 38, 39, 40, 41,
	42, 43, 44, 45, 46, 47, 51, 52, 53, 54,
	55, 56, 57, 58, 59, 60, 63, 64, 68, 69,
	72, 73, 74, 75, 81, 84, 85, 91, 92, 93,
	94, 95, 96,
}
var yyTok3 = []int{
	0,
}

//line yaccpar:1

/*	parser for yacc output	*/

var yyDebug = 0

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

const yyFlag = -1000

func yyTokname(c int) string {
	// 4 is TOKSTART above
	if c >= 4 && c-4 < len(yyToknames) {
		if yyToknames[c-4] != "" {
			return yyToknames[c-4]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func yylex1(lex yyLexer, lval *yySymType) int {
	c := 0
	char := lex.Lex(lval)
	if char <= 0 {
		c = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		c = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			c = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		c = yyTok3[i+0]
		if c == char {
			c = yyTok3[i+1]
			goto out
		}
	}

out:
	if c == 0 {
		c = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(c), uint(char))
	}
	return c
}

func yyParse(yylex yyLexer) int {
	var yyn int
	var yylval yySymType
	var yyVAL yySymType
	yyS := make([]yySymType, yyMaxDepth)

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yychar := -1
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yychar), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = yyPact[yystate]
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yychar < 0 {
		yychar = yylex1(yylex, &yylval)
	}
	yyn += yychar
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yychar { /* valid shift */
		yychar = -1
		yyVAL = yylval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = yyDef[yystate]
	if yyn == -2 {
		if yychar < 0 {
			yychar = yylex1(yylex, &yylval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && yyExca[xi+1] == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = yyExca[xi+0]
			if yyn < 0 || yyn == yychar {
				break
			}
		}
		yyn = yyExca[xi+1]
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error("syntax error")
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yychar))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = yyPact[yyS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yychar))
			}
			if yychar == yyEofCode {
				goto ret1
			}
			yychar = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= yyR2[yyn]
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = yyAct[yyg]
	} else {
		yystate = yyAct[yyj]
		if yyChk[yystate] != -yyn {
			yystate = yyAct[yyg]
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 1:
		//line cc.y:181
		{
			yylex.(*lexer).prog = &Prog{Decl: yyS[yypt-1].decls}
			return 0
		}
	case 2:
		//line cc.y:186
		{
			yylex.(*lexer).expr = yyS[yypt-1].expr
			return 0
		}
	case 3:
		//line cc.y:192
		{
			yyVAL.span = Span{}
			yyVAL.decls = nil
		}
	case 4:
		//line cc.y:197
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.decls = append(yyS[yypt-1].decls, yyS[yypt-0].decls...)
		}
	case 5:
		//line cc.y:204
		{
			yyVAL.span = yyS[yypt-0].span
			if len(yyS[yypt-0].exprs) == 1 {
				yyVAL.expr = yyS[yypt-0].exprs[0]
				break
			}
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: Comma, List: yyS[yypt-0].exprs}
		}
	case 6:
		//line cc.y:215
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: Name, Text: yyS[yypt-0].str}
		}
	case 7:
		//line cc.y:220
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: Number, Text: yyS[yypt-0].str}
		}
	case 8:
		//line cc.y:225
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: Number, Text: yyS[yypt-0].str}
		}
	case 9:
		//line cc.y:230
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: String, Texts: yyS[yypt-0].strs}
		}
	case 10:
		//line cc.y:235
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: Add, Left: yyS[yypt-2].expr, Right: yyS[yypt-0].expr}
		}
	case 11:
		//line cc.y:240
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: Sub, Left: yyS[yypt-2].expr, Right: yyS[yypt-0].expr}
		}
	case 12:
		//line cc.y:245
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: Mul, Left: yyS[yypt-2].expr, Right: yyS[yypt-0].expr}
		}
	case 13:
		//line cc.y:250
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: Div, Left: yyS[yypt-2].expr, Right: yyS[yypt-0].expr}
		}
	case 14:
		//line cc.y:255
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: Mod, Left: yyS[yypt-2].expr, Right: yyS[yypt-0].expr}
		}
	case 15:
		//line cc.y:260
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: Lsh, Left: yyS[yypt-2].expr, Right: yyS[yypt-0].expr}
		}
	case 16:
		//line cc.y:265
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: Rsh, Left: yyS[yypt-2].expr, Right: yyS[yypt-0].expr}
		}
	case 17:
		//line cc.y:270
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: Lt, Left: yyS[yypt-2].expr, Right: yyS[yypt-0].expr}
		}
	case 18:
		//line cc.y:275
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: Gt, Left: yyS[yypt-2].expr, Right: yyS[yypt-0].expr}
		}
	case 19:
		//line cc.y:280
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: LtEq, Left: yyS[yypt-2].expr, Right: yyS[yypt-0].expr}
		}
	case 20:
		//line cc.y:285
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: GtEq, Left: yyS[yypt-2].expr, Right: yyS[yypt-0].expr}
		}
	case 21:
		//line cc.y:290
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: EqEq, Left: yyS[yypt-2].expr, Right: yyS[yypt-0].expr}
		}
	case 22:
		//line cc.y:295
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: NotEq, Left: yyS[yypt-2].expr, Right: yyS[yypt-0].expr}
		}
	case 23:
		//line cc.y:300
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: And, Left: yyS[yypt-2].expr, Right: yyS[yypt-0].expr}
		}
	case 24:
		//line cc.y:305
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: Xor, Left: yyS[yypt-2].expr, Right: yyS[yypt-0].expr}
		}
	case 25:
		//line cc.y:310
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: Or, Left: yyS[yypt-2].expr, Right: yyS[yypt-0].expr}
		}
	case 26:
		//line cc.y:315
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: AndAnd, Left: yyS[yypt-2].expr, Right: yyS[yypt-0].expr}
		}
	case 27:
		//line cc.y:320
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: OrOr, Left: yyS[yypt-2].expr, Right: yyS[yypt-0].expr}
		}
	case 28:
		//line cc.y:325
		{
			yyVAL.span = span(yyS[yypt-4].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: Cond, List: []*Expr{yyS[yypt-4].expr, yyS[yypt-2].expr, yyS[yypt-0].expr}}
		}
	case 29:
		//line cc.y:330
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: Eq, Left: yyS[yypt-2].expr, Right: yyS[yypt-0].expr}
		}
	case 30:
		//line cc.y:335
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: AddEq, Left: yyS[yypt-2].expr, Right: yyS[yypt-0].expr}
		}
	case 31:
		//line cc.y:340
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: SubEq, Left: yyS[yypt-2].expr, Right: yyS[yypt-0].expr}
		}
	case 32:
		//line cc.y:345
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: MulEq, Left: yyS[yypt-2].expr, Right: yyS[yypt-0].expr}
		}
	case 33:
		//line cc.y:350
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: DivEq, Left: yyS[yypt-2].expr, Right: yyS[yypt-0].expr}
		}
	case 34:
		//line cc.y:355
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: ModEq, Left: yyS[yypt-2].expr, Right: yyS[yypt-0].expr}
		}
	case 35:
		//line cc.y:360
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: LshEq, Left: yyS[yypt-2].expr, Right: yyS[yypt-0].expr}
		}
	case 36:
		//line cc.y:365
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: RshEq, Left: yyS[yypt-2].expr, Right: yyS[yypt-0].expr}
		}
	case 37:
		//line cc.y:370
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: AndEq, Left: yyS[yypt-2].expr, Right: yyS[yypt-0].expr}
		}
	case 38:
		//line cc.y:375
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: XorEq, Left: yyS[yypt-2].expr, Right: yyS[yypt-0].expr}
		}
	case 39:
		//line cc.y:380
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: OrEq, Left: yyS[yypt-2].expr, Right: yyS[yypt-0].expr}
		}
	case 40:
		//line cc.y:385
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: Indir, Left: yyS[yypt-0].expr}
		}
	case 41:
		//line cc.y:390
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: Addr, Left: yyS[yypt-0].expr}
		}
	case 42:
		//line cc.y:395
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: Plus, Left: yyS[yypt-0].expr}
		}
	case 43:
		//line cc.y:400
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: Minus, Left: yyS[yypt-0].expr}
		}
	case 44:
		//line cc.y:405
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: Not, Left: yyS[yypt-0].expr}
		}
	case 45:
		//line cc.y:410
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: Twid, Left: yyS[yypt-0].expr}
		}
	case 46:
		//line cc.y:415
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: PreInc, Left: yyS[yypt-0].expr}
		}
	case 47:
		//line cc.y:420
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: PreDec, Left: yyS[yypt-0].expr}
		}
	case 48:
		//line cc.y:425
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: SizeofExpr, Left: yyS[yypt-0].expr}
		}
	case 49:
		//line cc.y:430
		{
			yyVAL.span = span(yyS[yypt-3].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: SizeofType, Type: yyS[yypt-1].typ}
		}
	case 50:
		//line cc.y:435
		{
			yyVAL.span = span(yyS[yypt-5].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: Offsetof, Type: yyS[yypt-3].typ, Left: yyS[yypt-1].expr}
		}
	case 51:
		//line cc.y:440
		{
			yyVAL.span = span(yyS[yypt-3].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: Cast, Type: yyS[yypt-2].typ, Left: yyS[yypt-0].expr}
		}
	case 52:
		//line cc.y:445
		{
			yyVAL.span = span(yyS[yypt-3].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: CastInit, Type: yyS[yypt-2].typ, Init: &Init{Span: yyVAL.span, Braced: yyS[yypt-0].inits}}
		}
	case 53:
		//line cc.y:450
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: Paren, Left: yyS[yypt-1].expr}
		}
	case 54:
		//line cc.y:455
		{
			yyVAL.span = span(yyS[yypt-3].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: Call, Left: yyS[yypt-3].expr, List: yyS[yypt-1].exprs}
		}
	case 55:
		//line cc.y:460
		{
			yyVAL.span = span(yyS[yypt-3].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: Index, Left: yyS[yypt-3].expr, Right: yyS[yypt-1].expr}
		}
	case 56:
		//line cc.y:465
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: PostInc, Left: yyS[yypt-1].expr}
		}
	case 57:
		//line cc.y:470
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: PostDec, Left: yyS[yypt-1].expr}
		}
	case 58:
		//line cc.y:475
		{
			yyVAL.span = span(yyS[yypt-5].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: VaArg, Left: yyS[yypt-3].expr, Type: yyS[yypt-1].typ}
		}
	case 59:
		//line cc.y:481
		{
			yyVAL.span = Span{}
			yyVAL.stmts = nil
		}
	case 60:
		//line cc.y:486
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.stmts = yyS[yypt-1].stmts
			for _, d := range yyS[yypt-0].decls {
				yyVAL.stmts = append(yyVAL.stmts, &Stmt{Span: yyVAL.span, Op: StmtDecl, Decl: d})
			}
		}
	case 61:
		//line cc.y:494
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.stmts = append(yyS[yypt-1].stmts, yyS[yypt-0].stmt)
		}
	case 62:
		//line cc.y:501
		{
			pushScope()
		}
	case 63:
		//line cc.y:505
		{
			yyVAL.span = span(yyS[yypt-3].span, yyS[yypt-0].span)
			popScope()
			yyVAL.stmt = &Stmt{Span: yyVAL.span, Op: Block, Block: yyS[yypt-1].stmts}
		}
	case 64:
		//line cc.y:513
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.label = &Label{Span: yyVAL.span, Op: Case, Expr: yyS[yypt-1].expr}
		}
	case 65:
		//line cc.y:518
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.label = &Label{Span: yyVAL.span, Op: Default}
		}
	case 66:
		//line cc.y:523
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.label = &Label{Span: yyVAL.span, Op: LabelName, Name: yyS[yypt-1].str}
		}
	case 67:
		//line cc.y:530
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.stmt = yyS[yypt-0].stmt
			yyVAL.stmt.Labels = yyS[yypt-1].labels
		}
	case 68:
		//line cc.y:538
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.stmt = &Stmt{Span: yyVAL.span, Op: Empty}
		}
	case 69:
		//line cc.y:543
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.stmt = yyS[yypt-0].stmt
		}
	case 70:
		//line cc.y:548
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.stmt = &Stmt{Span: yyVAL.span, Op: StmtExpr, Expr: yyS[yypt-1].expr}
		}
	case 71:
		//line cc.y:553
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.stmt = &Stmt{Span: yyVAL.span, Op: ARGBEGIN, Block: yyS[yypt-1].stmts}
		}
	case 72:
		//line cc.y:558
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.stmt = &Stmt{Span: yyVAL.span, Op: Break}
		}
	case 73:
		//line cc.y:563
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.stmt = &Stmt{Span: yyVAL.span, Op: Continue}
		}
	case 74:
		//line cc.y:568
		{
			yyVAL.span = span(yyS[yypt-6].span, yyS[yypt-0].span)
			yyVAL.stmt = &Stmt{Span: yyVAL.span, Op: Do, Body: yyS[yypt-5].stmt, Expr: yyS[yypt-2].expr}
		}
	case 75:
		//line cc.y:573
		{
			yyVAL.span = span(yyS[yypt-8].span, yyS[yypt-0].span)
			yyVAL.stmt = &Stmt{Span: yyVAL.span,
				Op:   For,
				Pre:  yyS[yypt-6].expr,
				Expr: yyS[yypt-4].expr,
				Post: yyS[yypt-2].expr,
				Body: yyS[yypt-0].stmt,
			}
		}
	case 76:
		//line cc.y:584
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.stmt = &Stmt{Span: yyVAL.span, Op: Goto, Text: yyS[yypt-1].str}
		}
	case 77:
		//line cc.y:589
		{
			yyVAL.span = span(yyS[yypt-4].span, yyS[yypt-0].span)
			yyVAL.stmt = &Stmt{Span: yyVAL.span, Op: If, Expr: yyS[yypt-2].expr, Body: yyS[yypt-0].stmt}
		}
	case 78:
		//line cc.y:594
		{
			yyVAL.span = span(yyS[yypt-6].span, yyS[yypt-0].span)
			yyVAL.stmt = &Stmt{Span: yyVAL.span, Op: If, Expr: yyS[yypt-4].expr, Body: yyS[yypt-2].stmt, Else: yyS[yypt-0].stmt}
		}
	case 79:
		//line cc.y:599
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.stmt = &Stmt{Span: yyVAL.span, Op: Return, Expr: yyS[yypt-1].expr}
		}
	case 80:
		//line cc.y:604
		{
			yyVAL.span = span(yyS[yypt-4].span, yyS[yypt-0].span)
			yyVAL.stmt = &Stmt{Span: yyVAL.span, Op: Switch, Expr: yyS[yypt-2].expr, Body: yyS[yypt-0].stmt}
		}
	case 81:
		//line cc.y:609
		{
			yyVAL.span = span(yyS[yypt-4].span, yyS[yypt-0].span)
			yyVAL.stmt = &Stmt{Span: yyVAL.span, Op: While, Expr: yyS[yypt-2].expr, Body: yyS[yypt-0].stmt}
		}
	case 82:
		//line cc.y:616
		{
			yyVAL.span = Span{}
			yyVAL.abdecor = func(t *Type) *Type { return t }
		}
	case 83:
		//line cc.y:621
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			_, q, _ := splitTypeWords(yyS[yypt-1].strs)
			abdecor := yyS[yypt-0].abdecor
			yyVAL.abdecor = func(t *Type) *Type {
				t = abdecor(t)
				t = &Type{Span: yyVAL.span, Kind: Ptr, Base: t, Qual: q}
				return t
			}
		}
	case 84:
		//line cc.y:632
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.abdecor = yyS[yypt-0].abdecor
		}
	case 85:
		//line cc.y:639
		{
			yyVAL.span = span(yyS[yypt-3].span, yyS[yypt-0].span)
			abdecor := yyS[yypt-3].abdecor
			decls := yyS[yypt-1].decls
			span := yyVAL.span
			yyVAL.abdecor = func(t *Type) *Type {
				t = abdecor(t)
				t = &Type{Span: span, Kind: Func, Base: t, Decls: decls}
				return t
			}
		}
	case 86:
		//line cc.y:651
		{
			yyVAL.span = span(yyS[yypt-3].span, yyS[yypt-0].span)
			abdecor := yyS[yypt-3].abdecor
			yyVAL.abdecor = func(t *Type) *Type {
				t = abdecor(t)
				// TODO: use expr
				return t
			}

		}
	case 87:
		//line cc.y:662
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.abdecor = yyS[yypt-1].abdecor
		}
	case 88:
		//line cc.y:670
		{
			yyVAL.span = yyS[yypt-0].span
			name := yyS[yypt-0].str
			yyVAL.decor = func(t *Type) (*Type, string) { return t, name }
		}
	case 89:
		//line cc.y:676
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			_, q, _ := splitTypeWords(yyS[yypt-1].strs)
			decor := yyS[yypt-0].decor
			span := yyVAL.span
			yyVAL.decor = func(t *Type) (*Type, string) {
				t, name := decor(t)
				t = &Type{Span: span, Kind: Ptr, Base: t, Qual: q}
				return t, name
			}
		}
	case 90:
		//line cc.y:688
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.decor = yyS[yypt-1].decor
		}
	case 91:
		//line cc.y:693
		{
			yyVAL.span = span(yyS[yypt-3].span, yyS[yypt-0].span)
			decor := yyS[yypt-3].decor
			decls := yyS[yypt-1].decls
			span := yyVAL.span
			yyVAL.decor = func(t *Type) (*Type, string) {
				t, name := decor(t)
				t = &Type{Span: span, Kind: Func, Base: t, Decls: decls}
				return t, name
			}
		}
	case 92:
		//line cc.y:705
		{
			yyVAL.span = span(yyS[yypt-3].span, yyS[yypt-0].span)
			decor := yyS[yypt-3].decor
			span := yyVAL.span
			yyVAL.decor = func(t *Type) (*Type, string) {
				t, name := decor(t)
				// TODO: use expr
				t = &Type{Span: span, Kind: Array, Base: t}
				return t, name
			}
		}
	case 93:
		//line cc.y:720
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.decl = &Decl{Span: yyVAL.span, Name: yyS[yypt-0].str}
		}
	case 94:
		//line cc.y:725
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.decl = &Decl{Span: yyVAL.span, Type: yyS[yypt-0].abdecor(yyS[yypt-1].typ)}
		}
	case 95:
		//line cc.y:730
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			typ, name := yyS[yypt-0].decor(yyS[yypt-1].typ)
			yyVAL.decl = &Decl{Span: yyVAL.span, Name: name, Type: typ}
		}
	case 96:
		//line cc.y:736
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.decl = &Decl{Span: yyVAL.span, Name: "..."}
		}
	case 97:
		//line cc.y:744
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.idec = idecor{yyS[yypt-0].decor, nil}
		}
	case 98:
		//line cc.y:749
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.idec = idecor{yyS[yypt-2].decor, yyS[yypt-0].init}
		}
	case 99:
		//line cc.y:757
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.str = yyS[yypt-0].str
		}
	case 100:
		//line cc.y:762
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.str = yyS[yypt-0].str
		}
	case 101:
		//line cc.y:767
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.str = yyS[yypt-0].str
		}
	case 102:
		//line cc.y:772
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.str = yyS[yypt-0].str
		}
	case 103:
		//line cc.y:777
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.str = yyS[yypt-0].str
		}
	case 104:
		//line cc.y:782
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.str = yyS[yypt-0].str
		}
	case 105:
		//line cc.y:790
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.str = yyS[yypt-0].str
		}
	case 106:
		//line cc.y:795
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.str = yyS[yypt-0].str
		}
	case 107:
		//line cc.y:803
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.str = yyS[yypt-0].str
		}
	case 108:
		//line cc.y:808
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.str = yyS[yypt-0].str
		}
	case 109:
		//line cc.y:813
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.str = yyS[yypt-0].str
		}
	case 110:
		//line cc.y:818
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.str = yyS[yypt-0].str
		}
	case 111:
		//line cc.y:823
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.str = yyS[yypt-0].str
		}
	case 112:
		//line cc.y:828
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.str = yyS[yypt-0].str
		}
	case 113:
		//line cc.y:833
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.str = yyS[yypt-0].str
		}
	case 114:
		//line cc.y:838
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.str = yyS[yypt-0].str
		}
	case 115:
		//line cc.y:843
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.str = yyS[yypt-0].str
		}
	case 116:
		//line cc.y:850
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.str = yyS[yypt-0].str
		}
	case 117:
		//line cc.y:855
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.str = yyS[yypt-0].str
		}
	case 118:
		//line cc.y:862
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.str = yyS[yypt-0].str
		}
	case 119:
		//line cc.y:867
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.str = yyS[yypt-0].str
		}
	case 120:
		//line cc.y:875
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.typ = yyS[yypt-0].typ
		}
	case 121:
		//line cc.y:888
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.tc.c, yyVAL.tc.q, yyVAL.tc.t = splitTypeWords(append(yyS[yypt-0].strs, "int"))
		}
	case 122:
		//line cc.y:893
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.tc.c, yyVAL.tc.q, _ = splitTypeWords(append(yyS[yypt-2].strs, yyS[yypt-0].strs...))
			yyVAL.tc.t = yyS[yypt-1].typ
		}
	case 123:
		//line cc.y:899
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyS[yypt-2].strs = append(yyS[yypt-2].strs, yyS[yypt-1].str)
			yyS[yypt-2].strs = append(yyS[yypt-2].strs, yyS[yypt-0].strs...)
			yyVAL.tc.c, yyVAL.tc.q, yyVAL.tc.t = splitTypeWords(yyS[yypt-2].strs)
		}
	case 124:
		//line cc.y:906
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.tc.c, yyVAL.tc.q, _ = splitTypeWords(yyS[yypt-0].strs)
			yyVAL.tc.t = yyS[yypt-1].typ
		}
	case 125:
		//line cc.y:912
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			var ts []string
			ts = append(ts, yyS[yypt-1].str)
			ts = append(ts, yyS[yypt-0].strs...)
			yyVAL.tc.c, yyVAL.tc.q, yyVAL.tc.t = splitTypeWords(ts)
		}
	case 126:
		//line cc.y:923
		{
			yyVAL.span = yyS[yypt-0].span
			if yyS[yypt-0].tc.c != 0 {
				yylex.(*lexer).Errorf("%v not allowed here", yyS[yypt-0].tc.c)
			}
			if yyS[yypt-0].tc.q != 0 {
				yylex.(*lexer).Errorf("%v ignored here (TODO)?", yyS[yypt-0].tc.q)
			}
			yyVAL.typ = yyS[yypt-0].tc.t
		}
	case 127:
		//line cc.y:936
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.typ = yyS[yypt-0].abdecor(yyS[yypt-1].typ)
		}
	case 128:
		//line cc.y:944
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			// TODO: use $1.q
			yyVAL.span = Span{}
			yyVAL.decls = nil
			for _, idec := range yyS[yypt-1].idecs {
				typ, name := idec.d(yyS[yypt-2].tc.t)
				if yyS[yypt-2].tc.c&Typedef != 0 {
					pushNamedType(name, typ)
				}
				yyVAL.decls = append(yyVAL.decls, &Decl{Span: yyVAL.span, Name: name, Type: typ, Storage: yyS[yypt-2].tc.c, Init: idec.i})
			}
			if yyS[yypt-1].idecs == nil {
				yyVAL.decls = append(yyVAL.decls, &Decl{Span: yyVAL.span, Name: "", Type: yyS[yypt-2].tc.t, Storage: yyS[yypt-2].tc.c})
			}
		}
	case 129:
		//line cc.y:964
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.decls = yyS[yypt-0].decls
		}
	case 130:
		//line cc.y:969
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.decls = []*Decl{yyS[yypt-0].decl}
		}
	case 131:
		//line cc.y:976
		{
			yyVAL.span = span(yyS[yypt-3].span, yyS[yypt-0].span)
			typ, name := yyS[yypt-2].decor(yyS[yypt-3].tc.t)
			yyVAL.decl = &Decl{Span: yyVAL.span, Name: name, Type: typ}
			if yyS[yypt-1].decls != nil {
				yylex.(*lexer).Errorf("cannot use pre-prototype definitions")
			}
			yyVAL.decl.Body = yyS[yypt-0].stmt
		}
	case 132:
		//line cc.y:988
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.str = yyS[yypt-0].str
		}
	case 133:
		//line cc.y:993
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.str = yyS[yypt-0].str
		}
	case 134:
		//line cc.y:1001
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.tk = Struct
		}
	case 135:
		//line cc.y:1006
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.tk = Union
		}
	case 136:
		//line cc.y:1013
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.decor = yyS[yypt-0].decor
		}
	case 137:
		//line cc.y:1018
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			name := yyS[yypt-2].str
			expr := yyS[yypt-0].expr
			yyVAL.decor = func(t *Type) (*Type, string) {
				t.Width = expr
				return t, name
			}
		}
	case 138:
		//line cc.y:1030
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.decls = nil
			for _, decor := range yyS[yypt-1].decors {
				typ, name := decor(yyS[yypt-2].typ)
				yyVAL.decls = append(yyVAL.decls, &Decl{Span: yyVAL.span, Name: name, Type: typ})
			}
			if yyS[yypt-1].decors == nil {
				yyVAL.decls = append(yyVAL.decls, &Decl{Span: yyVAL.span, Type: yyS[yypt-2].typ})
			}
		}
	case 139:
		//line cc.y:1044
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.typ = &Type{Span: yyVAL.span, Kind: yyS[yypt-1].tk, Tag: yyS[yypt-0].str}
		}
	case 140:
		//line cc.y:1049
		{
			yyVAL.span = span(yyS[yypt-4].span, yyS[yypt-0].span)
			yyVAL.typ = &Type{Span: yyVAL.span, Kind: yyS[yypt-4].tk, Tag: yyS[yypt-3].str, Decls: yyS[yypt-1].decls}
		}
	case 141:
		//line cc.y:1056
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.prefix = &Prefix{Span: yyVAL.span, Dot: yyS[yypt-0].str}
		}
	case 142:
		//line cc.y:1063
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: Arrow, Left: yyS[yypt-2].expr, Text: yyS[yypt-0].str}
		}
	case 143:
		//line cc.y:1068
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: Dot, Left: yyS[yypt-2].expr, Text: yyS[yypt-0].str}
		}
	case 144:
		//line cc.y:1076
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.typ = &Type{Span: yyVAL.span, Kind: Enum, Tag: yyS[yypt-0].str}
		}
	case 145:
		//line cc.y:1081
		{
			yyVAL.span = span(yyS[yypt-5].span, yyS[yypt-0].span)
			yyVAL.typ = &Type{Span: yyVAL.span, Kind: Enum, Tag: yyS[yypt-4].str, Decls: yyS[yypt-2].decls}
		}
	case 146:
		//line cc.y:1088
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.decl = &Decl{Span: yyVAL.span, Name: yyS[yypt-1].str, Init: &Init{Span: yyVAL.span, Expr: yyS[yypt-0].expr}}
		}
	case 147:
		//line cc.y:1095
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.expr = yyS[yypt-0].expr
		}
	case 148:
		//line cc.y:1103
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.init = &Init{Span: yyVAL.span, Expr: yyS[yypt-0].expr}
		}
	case 149:
		//line cc.y:1108
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.init = &Init{Span: yyVAL.span, Braced: yyS[yypt-0].inits}
		}
	case 150:
		//line cc.y:1115
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.inits = []*Init{}
		}
	case 151:
		//line cc.y:1120
		{
			yyVAL.span = span(yyS[yypt-3].span, yyS[yypt-0].span)
			yyVAL.inits = append(yyS[yypt-2].inits, yyS[yypt-1].init)
		}
	case 152:
		//line cc.y:1125
		{
			yyVAL.span = span(yyS[yypt-4].span, yyS[yypt-0].span)
			yyVAL.inits = append(yyS[yypt-3].inits, yyS[yypt-2].init)
		}
	case 153:
		//line cc.y:1131
		{
			yyVAL.span = Span{}
			yyVAL.inits = nil
		}
	case 154:
		//line cc.y:1136
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.inits = append(yyS[yypt-2].inits, yyS[yypt-1].init)
		}
	case 155:
		//line cc.y:1143
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.init = yyS[yypt-0].init
		}
	case 156:
		//line cc.y:1148
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.init = yyS[yypt-0].init
			yyVAL.init.Prefix = yyS[yypt-2].prefixes
		}
	case 157:
		//line cc.y:1156
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.prefix = &Prefix{Span: yyVAL.span, Index: yyS[yypt-1].expr}
		}
	case 158:
		//line cc.y:1162
		{
			yyVAL.span = Span{}
		}
	case 159:
		//line cc.y:1166
		{
			yyVAL.span = yyS[yypt-0].span
		}
	case 160:
		//line cc.y:1171
		{
			yyVAL.span = Span{}
		}
	case 161:
		//line cc.y:1175
		{
			yyVAL.span = yyS[yypt-0].span
		}
	case 162:
		//line cc.y:1184
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.prefixes = []*Prefix{yyS[yypt-0].prefix}
		}
	case 163:
		//line cc.y:1189
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.prefixes = append(yyS[yypt-1].prefixes, yyS[yypt-0].prefix)
		}
	case 164:
		//line cc.y:1195
		{
			yyVAL.span = Span{}
			yyVAL.str = ""
		}
	case 165:
		//line cc.y:1200
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.str = yyS[yypt-0].str
		}
	case 166:
		//line cc.y:1206
		{
			yyVAL.span = Span{}
			yyVAL.expr = nil
		}
	case 167:
		//line cc.y:1211
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.expr = yyS[yypt-0].expr
		}
	case 168:
		//line cc.y:1217
		{
			yyVAL.span = Span{}
			yyVAL.expr = nil
		}
	case 169:
		//line cc.y:1222
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.expr = yyS[yypt-0].expr
		}
	case 170:
		//line cc.y:1229
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.exprs = []*Expr{yyS[yypt-0].expr}
		}
	case 171:
		//line cc.y:1234
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.exprs = append(yyS[yypt-2].exprs, yyS[yypt-0].expr)
		}
	case 172:
		//line cc.y:1240
		{
			yyVAL.span = Span{}
			yyVAL.exprs = nil
		}
	case 173:
		//line cc.y:1245
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.exprs = yyS[yypt-0].exprs
		}
	case 174:
		//line cc.y:1251
		{
			yyVAL.span = Span{}
			yyVAL.decls = nil
		}
	case 175:
		//line cc.y:1256
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.decls = append(yyS[yypt-1].decls, yyS[yypt-0].decls...)
		}
	case 176:
		//line cc.y:1262
		{
			yyVAL.span = Span{}
			yyVAL.labels = nil
		}
	case 177:
		//line cc.y:1267
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.labels = append(yyS[yypt-1].labels, yyS[yypt-0].label)
		}
	case 178:
		//line cc.y:1274
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.decls = []*Decl{yyS[yypt-0].decl}
		}
	case 179:
		//line cc.y:1279
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.decls = append(yyS[yypt-2].decls, yyS[yypt-0].decl)
		}
	case 180:
		//line cc.y:1285
		{
			yyVAL.span = Span{}
			yyVAL.decls = nil
		}
	case 181:
		//line cc.y:1290
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.decls = yyS[yypt-0].decls
		}
	case 182:
		//line cc.y:1297
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.idecs = []idecor{yyS[yypt-0].idec}
		}
	case 183:
		//line cc.y:1302
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.idecs = append(yyS[yypt-2].idecs, yyS[yypt-0].idec)
		}
	case 184:
		//line cc.y:1308
		{
			yyVAL.span = Span{}
			yyVAL.idecs = nil
		}
	case 185:
		//line cc.y:1313
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.idecs = yyS[yypt-0].idecs
		}
	case 186:
		//line cc.y:1320
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.strs = []string{yyS[yypt-0].str}
		}
	case 187:
		//line cc.y:1325
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.strs = append(yyS[yypt-1].strs, yyS[yypt-0].str)
		}
	case 188:
		//line cc.y:1331
		{
			yyVAL.span = Span{}
			yyVAL.strs = nil
		}
	case 189:
		//line cc.y:1336
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.strs = yyS[yypt-0].strs
		}
	case 190:
		//line cc.y:1343
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.strs = []string{yyS[yypt-0].str}
		}
	case 191:
		//line cc.y:1348
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.strs = append(yyS[yypt-1].strs, yyS[yypt-0].str)
		}
	case 192:
		//line cc.y:1354
		{
			yyVAL.span = Span{}
			yyVAL.strs = nil
		}
	case 193:
		//line cc.y:1359
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.strs = yyS[yypt-0].strs
		}
	case 194:
		//line cc.y:1366
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.strs = []string{yyS[yypt-0].str}
		}
	case 195:
		//line cc.y:1371
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.strs = append(yyS[yypt-1].strs, yyS[yypt-0].str)
		}
	case 196:
		//line cc.y:1377
		{
			yyVAL.span = Span{}
			yyVAL.strs = nil
		}
	case 197:
		//line cc.y:1382
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.strs = yyS[yypt-0].strs
		}
	case 198:
		//line cc.y:1389
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.decors = nil
			yyVAL.decors = append(yyVAL.decors, yyS[yypt-0].decor)
		}
	case 199:
		//line cc.y:1395
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.decors = append(yyS[yypt-2].decors, yyS[yypt-0].decor)
		}
	case 200:
		//line cc.y:1401
		{
			yyVAL.span = Span{}
			yyVAL.decors = nil
		}
	case 201:
		//line cc.y:1406
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.decors = yyS[yypt-0].decors
		}
	case 202:
		//line cc.y:1413
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.decls = yyS[yypt-0].decls
		}
	case 203:
		//line cc.y:1418
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.decls = append(yyS[yypt-1].decls, yyS[yypt-0].decls...)
		}
	case 204:
		//line cc.y:1424
		{
			yyVAL.span = Span{}
			yyVAL.expr = nil
		}
	case 205:
		//line cc.y:1429
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.expr = yyS[yypt-0].expr
		}
	case 206:
		//line cc.y:1436
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.decls = []*Decl{yyS[yypt-0].decl}
		}
	case 207:
		//line cc.y:1441
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.decls = append(yyS[yypt-2].decls, yyS[yypt-0].decl)
		}
	case 208:
		//line cc.y:1448
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.strs = []string{yyS[yypt-0].str}
		}
	case 209:
		//line cc.y:1453
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.strs = append(yyS[yypt-1].strs, yyS[yypt-0].str)
		}
	}
	goto yystack /* stack new state and value */
}
