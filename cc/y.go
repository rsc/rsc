//line cc.y:34
package cc

import __yyfmt__ "fmt"

//line cc.y:34
type typeClass struct {
	c Storage
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
const tokAUTOLIB = 57348
const tokAuto = 57349
const tokBreak = 57350
const tokCase = 57351
const tokChar = 57352
const tokConst = 57353
const tokContinue = 57354
const tokDefault = 57355
const tokDo = 57356
const tokDotDotDot = 57357
const tokDouble = 57358
const tokEnum = 57359
const tokError = 57360
const tokExtern = 57361
const tokFloat = 57362
const tokFor = 57363
const tokGoto = 57364
const tokIf = 57365
const tokInline = 57366
const tokInt = 57367
const tokLitChar = 57368
const tokLong = 57369
const tokName = 57370
const tokNumber = 57371
const tokOffsetof = 57372
const tokRegister = 57373
const tokReturn = 57374
const tokShort = 57375
const tokSigned = 57376
const tokStatic = 57377
const tokStruct = 57378
const tokSwitch = 57379
const tokTypeName = 57380
const tokTypedef = 57381
const tokUnion = 57382
const tokUnsigned = 57383
const tokVaArg = 57384
const tokVoid = 57385
const tokVolatile = 57386
const tokWhile = 57387
const tokString = 57388
const tokShift = 57389
const tokElse = 57390
const tokAddEq = 57391
const tokSubEq = 57392
const tokMulEq = 57393
const tokDivEq = 57394
const tokModEq = 57395
const tokLshEq = 57396
const tokRshEq = 57397
const tokAndEq = 57398
const tokXorEq = 57399
const tokOrEq = 57400
const tokOrOr = 57401
const tokAndAnd = 57402
const tokEqEq = 57403
const tokNotEq = 57404
const tokLtEq = 57405
const tokGtEq = 57406
const tokLsh = 57407
const tokRsh = 57408
const tokCast = 57409
const tokSizeof = 57410
const tokUnary = 57411
const tokDec = 57412
const tokInc = 57413
const tokArrow = 57414
const startExpr = 57415
const startProg = 57416
const tokEOF = 57417

var yyToknames = []string{
	"tokARGBEGIN",
	"tokARGEND",
	"tokAUTOLIB",
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
	-1, 118,
	50, 98,
	99, 98,
	-2, 176,
	-1, 134,
	49, 167,
	-2, 141,
	-1, 138,
	49, 167,
	-2, 146,
	-1, 239,
	99, 202,
	-2, 166,
	-1, 268,
	63, 134,
	-2, 89,
}

const yyNprod = 212
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 1420

var yyAct = []int{

	7, 302, 112, 269, 28, 336, 287, 31, 228, 225,
	265, 196, 213, 98, 99, 100, 101, 102, 103, 104,
	105, 106, 241, 279, 337, 111, 219, 49, 5, 238,
	193, 109, 217, 223, 227, 123, 4, 131, 129, 365,
	127, 134, 138, 118, 363, 353, 32, 110, 348, 346,
	331, 330, 328, 288, 295, 187, 307, 291, 245, 59,
	120, 140, 141, 142, 143, 144, 145, 146, 147, 148,
	149, 150, 151, 152, 153, 154, 155, 156, 157, 158,
	128, 160, 161, 162, 163, 164, 165, 166, 167, 168,
	169, 170, 3, 2, 34, 189, 190, 188, 236, 174,
	175, 35, 294, 368, 159, 63, 64, 65, 362, 210,
	6, 257, 250, 96, 92, 184, 91, 173, 94, 93,
	95, 356, 258, 352, 180, 355, 354, 125, 189, 133,
	188, 339, 110, 189, 126, 188, 132, 284, 120, 176,
	177, 96, 92, 283, 91, 338, 94, 93, 95, 195,
	74, 72, 73, 68, 69, 70, 71, 66, 67, 61,
	62, 63, 64, 65, 198, 197, 253, 268, 199, 96,
	92, 128, 91, 335, 94, 93, 95, 137, 215, 207,
	66, 67, 61, 62, 63, 64, 65, 205, 203, 121,
	224, 226, 96, 92, 231, 91, 229, 94, 93, 95,
	122, 333, 172, 243, 233, 234, 207, 244, 179, 195,
	181, 224, 211, 208, 221, 178, 212, 115, 121, 31,
	133, 183, 235, 216, 248, 133, 239, 132, 232, 122,
	126, 114, 132, 256, 255, 298, 108, 282, 221, 259,
	208, 204, 210, 233, 247, 249, 251, 226, 342, 341,
	10, 266, 8, 9, 21, 290, 277, 202, 274, 289,
	271, 254, 206, 239, 60, 192, 23, 262, 201, 200,
	24, 280, 281, 209, 186, 293, 364, 116, 97, 344,
	221, 285, 300, 136, 299, 195, 57, 242, 286, 267,
	231, 306, 301, 137, 292, 185, 296, 226, 234, 248,
	305, 266, 297, 33, 270, 260, 308, 16, 17, 20,
	1, 37, 11, 313, 22, 194, 19, 18, 130, 58,
	332, 48, 329, 310, 334, 278, 340, 135, 139, 314,
	304, 311, 231, 246, 301, 276, 124, 117, 119, 345,
	171, 272, 273, 263, 264, 240, 237, 26, 218, 191,
	29, 182, 0, 0, 0, 0, 359, 360, 361, 358,
	347, 0, 0, 349, 350, 0, 367, 0, 0, 366,
	369, 0, 0, 0, 0, 0, 0, 357, 80, 81,
	82, 83, 84, 85, 86, 87, 88, 89, 90, 79,
	351, 78, 77, 76, 75, 74, 72, 73, 68, 69,
	70, 71, 66, 67, 61, 62, 63, 64, 65, 0,
	0, 0, 0, 0, 96, 92, 0, 91, 0, 94,
	93, 95, 80, 81, 82, 83, 84, 85, 86, 87,
	88, 89, 90, 79, 0, 78, 77, 76, 75, 74,
	72, 73, 68, 69, 70, 71, 66, 67, 61, 62,
	63, 64, 65, 0, 0, 0, 0, 0, 96, 92,
	309, 91, 0, 94, 93, 95, 80, 81, 82, 83,
	84, 85, 86, 87, 88, 89, 90, 79, 0, 78,
	77, 76, 75, 74, 72, 73, 68, 69, 70, 71,
	66, 67, 61, 62, 63, 64, 65, 0, 0, 0,
	0, 0, 96, 92, 0, 91, 275, 94, 93, 95,
	214, 80, 81, 82, 83, 84, 85, 86, 87, 88,
	89, 90, 79, 0, 78, 77, 76, 75, 74, 72,
	73, 68, 69, 70, 71, 66, 67, 61, 62, 63,
	64, 65, 0, 0, 0, 0, 0, 96, 92, 0,
	91, 0, 94, 93, 95, 80, 81, 82, 83, 84,
	85, 86, 87, 88, 89, 90, 79, 0, 78, 77,
	76, 75, 74, 72, 73, 68, 69, 70, 71, 66,
	67, 61, 62, 63, 64, 65, 0, 0, 0, 0,
	0, 96, 92, 0, 91, 0, 94, 93, 95, 52,
	0, 0, 39, 57, 0, 0, 0, 0, 46, 38,
	0, 113, 45, 0, 0, 0, 56, 41, 10, 42,
	8, 9, 21, 55, 0, 40, 43, 53, 50, 0,
	36, 54, 51, 44, 23, 47, 58, 0, 24, 77,
	76, 75, 74, 72, 73, 68, 69, 70, 71, 66,
	67, 61, 62, 63, 64, 65, 0, 0, 0, 0,
	13, 96, 92, 0, 91, 0, 94, 93, 95, 14,
	15, 12, 0, 0, 0, 16, 17, 20, 0, 0,
	0, 0, 22, 315, 19, 18, 0, 316, 325, 0,
	0, 317, 326, 318, 0, 0, 0, 0, 0, 0,
	319, 320, 321, 0, 0, 10, 0, 327, 9, 21,
	0, 322, 0, 0, 0, 0, 323, 0, 0, 0,
	0, 23, 0, 0, 324, 24, 0, 0, 230, 0,
	72, 73, 68, 69, 70, 71, 66, 67, 61, 62,
	63, 64, 65, 0, 0, 0, 0, 13, 96, 92,
	0, 91, 0, 94, 93, 95, 14, 15, 12, 0,
	0, 0, 16, 17, 20, 0, 0, 0, 0, 22,
	0, 19, 18, 0, 0, 27, 52, 0, 312, 39,
	57, 0, 0, 0, 0, 46, 38, 0, 30, 45,
	0, 0, 0, 56, 41, 0, 42, 0, 0, 0,
	55, 0, 40, 43, 53, 50, 0, 36, 54, 51,
	44, 52, 47, 58, 39, 57, 0, 0, 0, 0,
	46, 38, 0, 113, 45, 0, 0, 0, 56, 41,
	0, 42, 0, 0, 0, 55, 0, 40, 43, 53,
	50, 0, 36, 54, 51, 44, 52, 47, 58, 39,
	57, 0, 0, 0, 0, 46, 38, 0, 113, 45,
	0, 0, 0, 56, 41, 0, 42, 252, 0, 0,
	55, 0, 40, 43, 53, 50, 0, 36, 54, 51,
	44, 0, 47, 58, 0, 27, 52, 0, 0, 39,
	57, 0, 0, 0, 0, 46, 38, 0, 30, 45,
	0, 0, 303, 56, 41, 0, 42, 0, 0, 0,
	55, 0, 40, 43, 53, 50, 0, 36, 54, 51,
	44, 0, 47, 58, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 79, 261, 78, 77,
	76, 75, 74, 72, 73, 68, 69, 70, 71, 66,
	67, 61, 62, 63, 64, 65, 0, 0, 0, 0,
	0, 96, 92, 0, 91, 0, 94, 93, 95, 0,
	0, 0, 0, 0, 0, 0, 25, 76, 75, 74,
	72, 73, 68, 69, 70, 71, 66, 67, 61, 62,
	63, 64, 65, 0, 0, 0, 0, 0, 96, 92,
	0, 91, 0, 94, 93, 95, 75, 74, 72, 73,
	68, 69, 70, 71, 66, 67, 61, 62, 63, 64,
	65, 10, 0, 8, 9, 21, 96, 92, 0, 91,
	0, 94, 93, 95, 0, 0, 0, 23, 0, 0,
	0, 24, 0, 0, 209, 0, 0, 73, 68, 69,
	70, 71, 66, 67, 61, 62, 63, 64, 65, 0,
	0, 0, 0, 13, 96, 92, 0, 91, 0, 94,
	93, 95, 14, 15, 12, 0, 0, 0, 16, 17,
	20, 0, 280, 281, 0, 22, 0, 19, 18, 68,
	69, 70, 71, 66, 67, 61, 62, 63, 64, 65,
	10, 0, 8, 9, 21, 96, 92, 0, 91, 52,
	94, 93, 95, 57, 0, 0, 23, 0, 0, 0,
	24, 113, 0, 209, 0, 0, 56, 10, 0, 8,
	9, 21, 0, 55, 0, 0, 0, 53, 0, 0,
	0, 54, 13, 23, 0, 0, 58, 24, 0, 0,
	0, 14, 15, 12, 0, 0, 0, 16, 17, 20,
	0, 0, 0, 0, 22, 0, 19, 18, 0, 13,
	0, 0, 10, 0, 8, 9, 21, 0, 14, 15,
	12, 0, 0, 0, 16, 17, 20, 0, 23, 0,
	0, 22, 24, 19, 18, 61, 62, 63, 64, 65,
	0, 0, 0, 0, 0, 96, 92, 0, 91, 0,
	94, 93, 95, 0, 13, 0, 0, 0, 0, 0,
	0, 0, 0, 14, 15, 12, 0, 0, 0, 16,
	17, 20, 0, 0, 0, 0, 107, 52, 19, 18,
	39, 57, 0, 0, 0, 0, 46, 38, 0, 113,
	45, 0, 0, 0, 56, 41, 0, 42, 0, 0,
	0, 55, 0, 40, 43, 53, 50, 0, 36, 54,
	51, 44, 52, 47, 58, 39, 57, 0, 0, 230,
	222, 46, 38, 0, 113, 45, 0, 0, 0, 56,
	41, 0, 42, 220, 0, 0, 55, 0, 40, 43,
	53, 50, 0, 36, 54, 51, 44, 0, 47, 58,
	343, 0, 52, 0, 0, 39, 57, 0, 0, 0,
	0, 46, 38, 0, 113, 45, 0, 0, 0, 56,
	41, 0, 42, 0, 0, 0, 55, 0, 40, 43,
	53, 50, 0, 36, 54, 51, 44, 52, 47, 58,
	39, 57, 0, 0, 0, 0, 46, 38, 0, 113,
	45, 0, 0, 0, 56, 41, 0, 42, 0, 0,
	0, 55, 0, 40, 43, 53, 50, 0, 36, 54,
	51, 44, 52, 47, 58, 39, 57, 0, 0, 0,
	0, 46, 0, 0, 113, 45, 0, 0, 0, 56,
	41, 0, 42, 0, 0, 0, 55, 0, 40, 43,
	53, 0, 0, 0, 54, 0, 44, 0, 47, 58,
}
var yyPact = []int{

	-3, -1000, -1000, 1101, 879, -38, 214, 504, -1000, -1000,
	-1000, 232, 1101, 1101, 1101, 1101, 1101, 1101, 1101, 1101,
	1146, 146, 592, 141, -1000, -1000, -1000, 127, -1000, -1000,
	231, 110, 1340, 1102, 1375, -1000, -1000, 255, 255, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	1101, 1101, 1101, 1101, 1101, 1101, 1101, 1101, 1101, 1101,
	1101, 1101, 1101, 1101, 1101, 1101, 1101, 1101, 1101, 1101,
	1101, 1101, 1101, 1101, 1101, 1101, 1101, 1101, 1101, 1101,
	1101, 1101, 1101, -1000, -1000, 255, 255, -1000, 54, 54,
	54, 54, 54, 54, 54, 54, 54, 592, 1340, 124,
	117, 131, -1000, -1000, 1101, 267, 225, -44, 45, 215,
	-1000, 275, 110, -1000, 1102, 1375, -1000, -1000, 1102, -1000,
	1375, -1000, -1000, -1000, -1000, 220, -1000, -1000, -1000, 219,
	504, 26, 26, 54, 54, 54, 1118, 1118, 105, 105,
	105, 105, 977, 1018, 661, 82, 939, 911, 574, 194,
	504, 504, 504, 504, 504, 504, 504, 504, 504, 504,
	504, 97, 214, 152, -1000, -1000, 96, 212, 1074, -1000,
	154, 275, 126, 131, 460, 87, -1000, -1000, 1265, 1101,
	1074, 1230, 110, 110, 275, -1000, 7, -1000, -1000, -1000,
	1340, 259, 1101, -1000, -1000, 224, 1101, 54, -1000, -40,
	1101, 131, 1265, 21, 1340, -1000, 769, 75, 211, -1000,
	-1000, 32, -1000, 150, 504, -1000, 504, -1000, -1000, -1000,
	-1000, 110, -1000, 45, 40, -1000, -1000, 839, -1000, 139,
	210, -1000, 207, 874, 415, -1000, 995, 148, 154, 52,
	-1000, 46, -1000, -1000, 1265, 154, 40, 275, 32, -1000,
	-1000, -1000, -1000, -46, 209, -1000, 40, 192, -1000, -1000,
	-41, 259, -1000, -1000, 1101, -1000, 4, -1000, 184, -1000,
	255, 1101, -1000, -1000, -1000, -1000, 32, 804, -1000, 139,
	1101, -1000, -1000, 504, -1000, -42, 1074, -1000, -1000, -1000,
	371, -1000, -1000, -1000, 679, -1000, 504, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -47, -1000, -48, -49, -1000, 111,
	255, 83, 1101, 55, 41, 1101, 186, 185, -1000, 1305,
	-1000, -1000, 234, 1101, -50, 1101, -51, -1000, 1101, 1101,
	327, -1000, -1000, -1000, 33, -54, -1000, 35, -1000, 34,
	30, -1000, 1101, 1101, -1000, -1000, -1000, 17, -55, 228,
	-1000, -1000, -60, 1101, -1000, -1000, 12, -1000, -1000, -1000,
}
var yyPgo = []int{

	0, 12, 351, 26, 350, 22, 4, 349, 348, 32,
	36, 347, 29, 346, 345, 11, 10, 344, 343, 0,
	33, 24, 5, 342, 341, 110, 340, 35, 338, 337,
	9, 335, 34, 333, 331, 330, 23, 325, 323, 8,
	1, 6, 321, 27, 94, 101, 37, 3, 289, 46,
	40, 318, 38, 315, 30, 312, 2, 311, 31, 25,
	303, 310, 305, 304, 296,
}
var yyR1 = []int{

	0, 61, 61, 10, 10, 10, 21, 19, 19, 19,
	19, 19, 19, 19, 19, 19, 19, 19, 19, 19,
	19, 19, 19, 19, 19, 19, 19, 19, 19, 19,
	19, 19, 19, 19, 19, 19, 19, 19, 19, 19,
	19, 19, 19, 19, 19, 19, 19, 19, 19, 19,
	19, 19, 19, 19, 19, 19, 19, 19, 19, 19,
	41, 41, 41, 62, 39, 34, 34, 34, 40, 38,
	38, 38, 38, 38, 38, 38, 38, 38, 38, 38,
	38, 38, 38, 1, 1, 1, 2, 2, 2, 15,
	15, 15, 15, 15, 3, 3, 3, 3, 27, 27,
	42, 42, 42, 42, 42, 42, 43, 43, 44, 44,
	44, 44, 44, 44, 44, 44, 44, 45, 45, 46,
	46, 60, 56, 56, 56, 56, 56, 59, 58, 6,
	11, 11, 11, 4, 47, 47, 57, 57, 16, 16,
	12, 60, 60, 36, 19, 19, 60, 60, 5, 23,
	30, 30, 32, 32, 32, 33, 33, 31, 31, 36,
	64, 64, 63, 63, 37, 37, 48, 48, 22, 22,
	20, 20, 25, 25, 26, 26, 7, 7, 35, 35,
	8, 8, 9, 9, 28, 28, 29, 29, 53, 53,
	54, 54, 49, 49, 50, 50, 51, 51, 52, 52,
	17, 17, 18, 18, 13, 13, 24, 24, 14, 14,
	55, 55,
}
var yyR2 = []int{

	0, 3, 3, 0, 2, 5, 1, 1, 1, 1,
	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 5,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 2, 2, 2, 2, 2, 2, 2, 2, 2,
	4, 6, 4, 4, 3, 4, 4, 2, 2, 6,
	0, 2, 2, 0, 4, 3, 2, 2, 2, 1,
	1, 2, 3, 2, 2, 7, 9, 3, 5, 7,
	3, 5, 5, 0, 3, 1, 4, 4, 3, 1,
	3, 3, 4, 4, 1, 2, 2, 1, 1, 3,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 3, 3, 2, 2, 1, 2, 3,
	1, 1, 5, 4, 1, 1, 1, 1, 1, 3,
	3, 2, 5, 2, 3, 3, 2, 6, 2, 2,
	1, 1, 2, 4, 5, 0, 3, 1, 3, 3,
	0, 1, 0, 1, 1, 2, 0, 1, 0, 1,
	0, 1, 1, 3, 0, 1, 0, 2, 0, 2,
	1, 3, 0, 1, 1, 3, 0, 1, 1, 2,
	0, 1, 1, 2, 0, 1, 1, 2, 0, 1,
	1, 3, 0, 1, 1, 2, 0, 1, 1, 3,
	1, 2,
}
var yyChk = []int{

	-1000, -61, 96, 95, -10, -21, -25, -19, 28, 29,
	26, -55, 79, 68, 77, 78, 83, 84, 93, 92,
	85, 30, 90, 42, 46, 97, -11, 6, -6, -4,
	19, -56, -49, -60, -44, -45, 38, -57, 17, 10,
	33, 25, 27, 34, 41, 20, 16, 43, -42, -43,
	36, 40, 7, 35, 39, 31, 24, 11, 44, 97,
	50, 77, 78, 79, 80, 81, 75, 76, 71, 72,
	73, 74, 69, 70, 68, 67, 66, 65, 64, 62,
	51, 52, 53, 54, 55, 56, 57, 58, 59, 60,
	61, 90, 88, 93, 92, 94, 87, 46, -19, -19,
	-19, -19, -19, -19, -19, -19, -19, 90, 90, -58,
	-21, -59, -56, 19, 90, 90, 46, -29, -15, -28,
	28, 79, 90, -27, -60, -44, -45, -50, -49, -52,
	-51, -46, -45, -44, -47, -48, 28, 38, -47, -48,
	-19, -19, -19, -19, -19, -19, -19, -19, -19, -19,
	-19, -19, -19, -19, -19, -19, -19, -19, -19, -21,
	-19, -19, -19, -19, -19, -19, -19, -19, -19, -19,
	-19, -26, -25, -21, -47, -47, -58, -58, 91, 91,
	-1, 79, -2, 90, -19, 28, 49, 99, 90, 88,
	51, -7, 50, -54, -53, -43, -15, -50, -52, -46,
	49, 49, 63, 91, 89, 91, 50, -19, -32, 49,
	88, -54, 90, -1, 50, 91, -10, -9, -8, -3,
	28, -59, 15, -20, -19, -30, -19, -32, -39, -6,
	49, -56, -27, -15, -15, -43, 91, -13, -12, -59,
	-14, -5, 28, -19, -19, 98, -33, -20, -1, -9,
	91, -58, 98, 91, 50, -1, -15, 79, 90, 89,
	-62, 98, -12, -18, -17, -16, -15, -48, 28, -47,
	-63, 50, -24, -23, 51, 91, -31, -30, -37, -36,
	87, 88, 89, 91, 91, -3, -54, -41, 99, 50,
	63, 98, -5, -19, 98, 50, -64, -36, 51, -47,
	-19, -6, -40, 98, -35, -16, -19, 98, -30, 89,
	-38, -34, 99, -39, -21, 4, 8, 12, 14, 21,
	22, 23, 32, 37, 45, 9, 13, 28, 99, -41,
	99, 99, -40, 90, -47, 90, -22, -21, 90, 90,
	-19, 63, 63, 5, 45, -22, 99, -21, 99, -21,
	-21, 63, 90, 99, 91, 91, 91, -21, -22, -40,
	-40, -40, 91, 99, 48, 99, -22, -40, 91, -40,
}
var yyDef = []int{

	0, -2, 3, 0, 0, 0, 6, 172, 7, 8,
	9, 10, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 210, 1, 4, 0, 130, 131,
	102, 186, 122, 194, 198, 192, 121, 166, 166, 108,
	109, 110, 111, 112, 113, 114, 115, 116, 117, 118,
	136, 137, 100, 101, 103, 104, 105, 106, 107, 2,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 174, 0, 57, 58, 0, 0, 211, 41, 42,
	43, 44, 45, 46, 47, 48, 49, 0, 0, 0,
	0, 83, 127, 102, 0, 0, 0, 0, -2, 187,
	89, 190, 0, 184, 194, 198, 193, 125, 195, 126,
	199, 196, 119, 120, -2, 0, 134, 135, -2, 0,
	173, 11, 12, 13, 14, 15, 16, 17, 18, 19,
	20, 21, 22, 23, 24, 25, 26, 27, 28, 0,
	30, 31, 32, 33, 34, 35, 36, 37, 38, 39,
	40, 0, 175, 0, 144, 145, 0, 0, 0, 54,
	128, 190, 85, 83, 0, 0, 3, 129, 182, 170,
	0, 0, 0, 0, 191, 188, 0, 123, 124, 197,
	0, 0, 0, 55, 56, 50, 0, 52, 53, 155,
	170, 83, 182, 0, 0, 5, 0, 0, 183, 180,
	94, 83, 97, 0, 171, 99, 150, 151, 133, 177,
	63, 186, 185, 98, 90, 189, 91, 0, 204, -2,
	162, 208, 206, 29, 0, 152, 0, 0, 84, 0,
	88, 0, 132, 92, 0, 95, 96, 190, 83, 93,
	60, 142, 205, 0, 203, 200, 138, 0, -2, 167,
	0, 163, 148, 207, 0, 51, 0, 157, 160, 164,
	0, 0, 87, 86, 59, 181, 83, 178, 140, 166,
	0, 147, 209, 149, 153, 156, 0, 165, 161, 143,
	0, 61, 62, 64, 0, 201, 139, 154, 158, 159,
	68, 179, 69, 70, 0, 60, 0, 0, 178, 0,
	0, 0, 168, 0, 0, 0, 0, 7, 71, 178,
	73, 74, 0, 168, 0, 0, 0, 169, 0, 0,
	0, 66, 67, 72, 0, 0, 77, 0, 80, 0,
	0, 65, 0, 168, 178, 178, 178, 0, 0, 78,
	81, 82, 0, 168, 178, 75, 0, 79, 178, 76,
}
var yyTok1 = []int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 83, 3, 3, 3, 81, 68, 3,
	90, 91, 79, 77, 50, 78, 87, 80, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 63, 99,
	71, 51, 72, 62, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 88, 3, 89, 67, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 49, 66, 98, 84,
}
var yyTok2 = []int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 36, 37, 38, 39, 40, 41,
	42, 43, 44, 45, 46, 47, 48, 52, 53, 54,
	55, 56, 57, 58, 59, 60, 61, 64, 65, 69,
	70, 73, 74, 75, 76, 82, 85, 86, 92, 93,
	94, 95, 96, 97,
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
		//line cc.y:182
		{
			yylex.(*lexer).prog = &Prog{Decl: yyS[yypt-1].decls}
			return 0
		}
	case 2:
		//line cc.y:187
		{
			yylex.(*lexer).expr = yyS[yypt-1].expr
			return 0
		}
	case 3:
		//line cc.y:193
		{
			yyVAL.span = Span{}
			yyVAL.decls = nil
		}
	case 4:
		//line cc.y:198
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.decls = append(yyS[yypt-1].decls, yyS[yypt-0].decls...)
		}
	case 5:
		//line cc.y:203
		{
		}
	case 6:
		//line cc.y:208
		{
			yyVAL.span = yyS[yypt-0].span
			if len(yyS[yypt-0].exprs) == 1 {
				yyVAL.expr = yyS[yypt-0].exprs[0]
				break
			}
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: Comma, List: yyS[yypt-0].exprs}
		}
	case 7:
		//line cc.y:219
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: Name, Text: yyS[yypt-0].str}
		}
	case 8:
		//line cc.y:224
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: Number, Text: yyS[yypt-0].str}
		}
	case 9:
		//line cc.y:229
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: Number, Text: yyS[yypt-0].str}
		}
	case 10:
		//line cc.y:234
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: String, Texts: yyS[yypt-0].strs}
		}
	case 11:
		//line cc.y:239
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: Add, Left: yyS[yypt-2].expr, Right: yyS[yypt-0].expr}
		}
	case 12:
		//line cc.y:244
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: Sub, Left: yyS[yypt-2].expr, Right: yyS[yypt-0].expr}
		}
	case 13:
		//line cc.y:249
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: Mul, Left: yyS[yypt-2].expr, Right: yyS[yypt-0].expr}
		}
	case 14:
		//line cc.y:254
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: Div, Left: yyS[yypt-2].expr, Right: yyS[yypt-0].expr}
		}
	case 15:
		//line cc.y:259
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: Mod, Left: yyS[yypt-2].expr, Right: yyS[yypt-0].expr}
		}
	case 16:
		//line cc.y:264
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: Lsh, Left: yyS[yypt-2].expr, Right: yyS[yypt-0].expr}
		}
	case 17:
		//line cc.y:269
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: Rsh, Left: yyS[yypt-2].expr, Right: yyS[yypt-0].expr}
		}
	case 18:
		//line cc.y:274
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: Lt, Left: yyS[yypt-2].expr, Right: yyS[yypt-0].expr}
		}
	case 19:
		//line cc.y:279
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: Gt, Left: yyS[yypt-2].expr, Right: yyS[yypt-0].expr}
		}
	case 20:
		//line cc.y:284
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: LtEq, Left: yyS[yypt-2].expr, Right: yyS[yypt-0].expr}
		}
	case 21:
		//line cc.y:289
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: GtEq, Left: yyS[yypt-2].expr, Right: yyS[yypt-0].expr}
		}
	case 22:
		//line cc.y:294
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: EqEq, Left: yyS[yypt-2].expr, Right: yyS[yypt-0].expr}
		}
	case 23:
		//line cc.y:299
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: NotEq, Left: yyS[yypt-2].expr, Right: yyS[yypt-0].expr}
		}
	case 24:
		//line cc.y:304
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: And, Left: yyS[yypt-2].expr, Right: yyS[yypt-0].expr}
		}
	case 25:
		//line cc.y:309
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: Xor, Left: yyS[yypt-2].expr, Right: yyS[yypt-0].expr}
		}
	case 26:
		//line cc.y:314
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: Or, Left: yyS[yypt-2].expr, Right: yyS[yypt-0].expr}
		}
	case 27:
		//line cc.y:319
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: AndAnd, Left: yyS[yypt-2].expr, Right: yyS[yypt-0].expr}
		}
	case 28:
		//line cc.y:324
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: OrOr, Left: yyS[yypt-2].expr, Right: yyS[yypt-0].expr}
		}
	case 29:
		//line cc.y:329
		{
			yyVAL.span = span(yyS[yypt-4].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: Cond, List: []*Expr{yyS[yypt-4].expr, yyS[yypt-2].expr, yyS[yypt-0].expr}}
		}
	case 30:
		//line cc.y:334
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: Eq, Left: yyS[yypt-2].expr, Right: yyS[yypt-0].expr}
		}
	case 31:
		//line cc.y:339
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: AddEq, Left: yyS[yypt-2].expr, Right: yyS[yypt-0].expr}
		}
	case 32:
		//line cc.y:344
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: SubEq, Left: yyS[yypt-2].expr, Right: yyS[yypt-0].expr}
		}
	case 33:
		//line cc.y:349
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: MulEq, Left: yyS[yypt-2].expr, Right: yyS[yypt-0].expr}
		}
	case 34:
		//line cc.y:354
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: DivEq, Left: yyS[yypt-2].expr, Right: yyS[yypt-0].expr}
		}
	case 35:
		//line cc.y:359
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: ModEq, Left: yyS[yypt-2].expr, Right: yyS[yypt-0].expr}
		}
	case 36:
		//line cc.y:364
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: LshEq, Left: yyS[yypt-2].expr, Right: yyS[yypt-0].expr}
		}
	case 37:
		//line cc.y:369
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: RshEq, Left: yyS[yypt-2].expr, Right: yyS[yypt-0].expr}
		}
	case 38:
		//line cc.y:374
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: AndEq, Left: yyS[yypt-2].expr, Right: yyS[yypt-0].expr}
		}
	case 39:
		//line cc.y:379
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: XorEq, Left: yyS[yypt-2].expr, Right: yyS[yypt-0].expr}
		}
	case 40:
		//line cc.y:384
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: OrEq, Left: yyS[yypt-2].expr, Right: yyS[yypt-0].expr}
		}
	case 41:
		//line cc.y:389
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: Indir, Left: yyS[yypt-0].expr}
		}
	case 42:
		//line cc.y:394
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: Addr, Left: yyS[yypt-0].expr}
		}
	case 43:
		//line cc.y:399
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: Plus, Left: yyS[yypt-0].expr}
		}
	case 44:
		//line cc.y:404
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: Minus, Left: yyS[yypt-0].expr}
		}
	case 45:
		//line cc.y:409
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: Not, Left: yyS[yypt-0].expr}
		}
	case 46:
		//line cc.y:414
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: Twid, Left: yyS[yypt-0].expr}
		}
	case 47:
		//line cc.y:419
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: PreInc, Left: yyS[yypt-0].expr}
		}
	case 48:
		//line cc.y:424
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: PreDec, Left: yyS[yypt-0].expr}
		}
	case 49:
		//line cc.y:429
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: SizeofExpr, Left: yyS[yypt-0].expr}
		}
	case 50:
		//line cc.y:434
		{
			yyVAL.span = span(yyS[yypt-3].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: SizeofType, Type: yyS[yypt-1].typ}
		}
	case 51:
		//line cc.y:439
		{
			yyVAL.span = span(yyS[yypt-5].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: Offsetof, Type: yyS[yypt-3].typ, Left: yyS[yypt-1].expr}
		}
	case 52:
		//line cc.y:444
		{
			yyVAL.span = span(yyS[yypt-3].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: Cast, Type: yyS[yypt-2].typ, Left: yyS[yypt-0].expr}
		}
	case 53:
		//line cc.y:449
		{
			yyVAL.span = span(yyS[yypt-3].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: CastInit, Type: yyS[yypt-2].typ, Init: &Init{Span: yyVAL.span, Braced: yyS[yypt-0].inits}}
		}
	case 54:
		//line cc.y:454
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: Paren, Left: yyS[yypt-1].expr}
		}
	case 55:
		//line cc.y:459
		{
			yyVAL.span = span(yyS[yypt-3].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: Call, Left: yyS[yypt-3].expr, List: yyS[yypt-1].exprs}
		}
	case 56:
		//line cc.y:464
		{
			yyVAL.span = span(yyS[yypt-3].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: Index, Left: yyS[yypt-3].expr, Right: yyS[yypt-1].expr}
		}
	case 57:
		//line cc.y:469
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: PostInc, Left: yyS[yypt-1].expr}
		}
	case 58:
		//line cc.y:474
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: PostDec, Left: yyS[yypt-1].expr}
		}
	case 59:
		//line cc.y:479
		{
			yyVAL.span = span(yyS[yypt-5].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: VaArg, Left: yyS[yypt-3].expr, Type: yyS[yypt-1].typ}
		}
	case 60:
		//line cc.y:485
		{
			yyVAL.span = Span{}
			yyVAL.stmts = nil
		}
	case 61:
		//line cc.y:490
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.stmts = yyS[yypt-1].stmts
			for _, d := range yyS[yypt-0].decls {
				yyVAL.stmts = append(yyVAL.stmts, &Stmt{Span: yyVAL.span, Op: StmtDecl, Decl: d})
			}
		}
	case 62:
		//line cc.y:498
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.stmts = append(yyS[yypt-1].stmts, yyS[yypt-0].stmt)
		}
	case 63:
		//line cc.y:505
		{
			yylex.(*lexer).pushScope()
		}
	case 64:
		//line cc.y:509
		{
			yyVAL.span = span(yyS[yypt-3].span, yyS[yypt-0].span)
			yylex.(*lexer).popScope()
			yyVAL.stmt = &Stmt{Span: yyVAL.span, Op: Block, Block: yyS[yypt-1].stmts}
		}
	case 65:
		//line cc.y:517
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.label = &Label{Span: yyVAL.span, Op: Case, Expr: yyS[yypt-1].expr}
		}
	case 66:
		//line cc.y:522
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.label = &Label{Span: yyVAL.span, Op: Default}
		}
	case 67:
		//line cc.y:527
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.label = &Label{Span: yyVAL.span, Op: LabelName, Name: yyS[yypt-1].str}
		}
	case 68:
		//line cc.y:534
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.stmt = yyS[yypt-0].stmt
			yyVAL.stmt.Labels = yyS[yypt-1].labels
		}
	case 69:
		//line cc.y:542
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.stmt = &Stmt{Span: yyVAL.span, Op: Empty}
		}
	case 70:
		//line cc.y:547
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.stmt = yyS[yypt-0].stmt
		}
	case 71:
		//line cc.y:552
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.stmt = &Stmt{Span: yyVAL.span, Op: StmtExpr, Expr: yyS[yypt-1].expr}
		}
	case 72:
		//line cc.y:557
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.stmt = &Stmt{Span: yyVAL.span, Op: ARGBEGIN, Block: yyS[yypt-1].stmts}
		}
	case 73:
		//line cc.y:562
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.stmt = &Stmt{Span: yyVAL.span, Op: Break}
		}
	case 74:
		//line cc.y:567
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.stmt = &Stmt{Span: yyVAL.span, Op: Continue}
		}
	case 75:
		//line cc.y:572
		{
			yyVAL.span = span(yyS[yypt-6].span, yyS[yypt-0].span)
			yyVAL.stmt = &Stmt{Span: yyVAL.span, Op: Do, Body: yyS[yypt-5].stmt, Expr: yyS[yypt-2].expr}
		}
	case 76:
		//line cc.y:577
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
	case 77:
		//line cc.y:588
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.stmt = &Stmt{Span: yyVAL.span, Op: Goto, Text: yyS[yypt-1].str}
		}
	case 78:
		//line cc.y:593
		{
			yyVAL.span = span(yyS[yypt-4].span, yyS[yypt-0].span)
			yyVAL.stmt = &Stmt{Span: yyVAL.span, Op: If, Expr: yyS[yypt-2].expr, Body: yyS[yypt-0].stmt}
		}
	case 79:
		//line cc.y:598
		{
			yyVAL.span = span(yyS[yypt-6].span, yyS[yypt-0].span)
			yyVAL.stmt = &Stmt{Span: yyVAL.span, Op: If, Expr: yyS[yypt-4].expr, Body: yyS[yypt-2].stmt, Else: yyS[yypt-0].stmt}
		}
	case 80:
		//line cc.y:603
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.stmt = &Stmt{Span: yyVAL.span, Op: Return, Expr: yyS[yypt-1].expr}
		}
	case 81:
		//line cc.y:608
		{
			yyVAL.span = span(yyS[yypt-4].span, yyS[yypt-0].span)
			yyVAL.stmt = &Stmt{Span: yyVAL.span, Op: Switch, Expr: yyS[yypt-2].expr, Body: yyS[yypt-0].stmt}
		}
	case 82:
		//line cc.y:613
		{
			yyVAL.span = span(yyS[yypt-4].span, yyS[yypt-0].span)
			yyVAL.stmt = &Stmt{Span: yyVAL.span, Op: While, Expr: yyS[yypt-2].expr, Body: yyS[yypt-0].stmt}
		}
	case 83:
		//line cc.y:620
		{
			yyVAL.span = Span{}
			yyVAL.abdecor = func(t *Type) *Type { return t }
		}
	case 84:
		//line cc.y:625
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			_, q, _ := splitTypeWords(yyS[yypt-1].strs)
			abdecor := yyS[yypt-0].abdecor
			yyVAL.abdecor = func(t *Type) *Type {
				return abdecor(&Type{Span: yyVAL.span, Kind: Ptr, Base: t, Qual: q})
			}
		}
	case 85:
		//line cc.y:634
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.abdecor = yyS[yypt-0].abdecor
		}
	case 86:
		//line cc.y:641
		{
			yyVAL.span = span(yyS[yypt-3].span, yyS[yypt-0].span)
			abdecor := yyS[yypt-3].abdecor
			decls := yyS[yypt-1].decls
			span := yyVAL.span
			yyVAL.abdecor = func(t *Type) *Type {
				return abdecor(&Type{Span: span, Kind: Func, Base: t, Decls: decls})
			}
		}
	case 87:
		//line cc.y:651
		{
			yyVAL.span = span(yyS[yypt-3].span, yyS[yypt-0].span)
			abdecor := yyS[yypt-3].abdecor
			span := yyVAL.span
			expr := yyS[yypt-1].expr
			yyVAL.abdecor = func(t *Type) *Type {
				return abdecor(&Type{Span: span, Kind: Array, Base: t, Width: expr})
			}

		}
	case 88:
		//line cc.y:662
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.abdecor = yyS[yypt-1].abdecor
		}
	case 89:
		//line cc.y:670
		{
			yyVAL.span = yyS[yypt-0].span
			name := yyS[yypt-0].str
			yyVAL.decor = func(t *Type) (*Type, string) { return t, name }
		}
	case 90:
		//line cc.y:676
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			_, q, _ := splitTypeWords(yyS[yypt-1].strs)
			decor := yyS[yypt-0].decor
			span := yyVAL.span
			yyVAL.decor = func(t *Type) (*Type, string) {
				return decor(&Type{Span: span, Kind: Ptr, Base: t, Qual: q})
			}
		}
	case 91:
		//line cc.y:686
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.decor = yyS[yypt-1].decor
		}
	case 92:
		//line cc.y:691
		{
			yyVAL.span = span(yyS[yypt-3].span, yyS[yypt-0].span)
			decor := yyS[yypt-3].decor
			decls := yyS[yypt-1].decls
			span := yyVAL.span
			yyVAL.decor = func(t *Type) (*Type, string) {
				return decor(&Type{Span: span, Kind: Func, Base: t, Decls: decls})
			}
		}
	case 93:
		//line cc.y:701
		{
			yyVAL.span = span(yyS[yypt-3].span, yyS[yypt-0].span)
			decor := yyS[yypt-3].decor
			span := yyVAL.span
			expr := yyS[yypt-1].expr
			yyVAL.decor = func(t *Type) (*Type, string) {
				return decor(&Type{Span: span, Kind: Array, Base: t, Width: expr})
			}
		}
	case 94:
		//line cc.y:714
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.decl = &Decl{Span: yyVAL.span, Name: yyS[yypt-0].str}
		}
	case 95:
		//line cc.y:719
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.decl = &Decl{Span: yyVAL.span, Type: yyS[yypt-0].abdecor(yyS[yypt-1].typ)}
		}
	case 96:
		//line cc.y:724
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			typ, name := yyS[yypt-0].decor(yyS[yypt-1].typ)
			yyVAL.decl = &Decl{Span: yyVAL.span, Name: name, Type: typ}
		}
	case 97:
		//line cc.y:730
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.decl = &Decl{Span: yyVAL.span, Name: "..."}
		}
	case 98:
		//line cc.y:738
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.idec = idecor{yyS[yypt-0].decor, nil}
		}
	case 99:
		//line cc.y:743
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.idec = idecor{yyS[yypt-2].decor, yyS[yypt-0].init}
		}
	case 100:
		//line cc.y:751
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.str = yyS[yypt-0].str
		}
	case 101:
		//line cc.y:756
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.str = yyS[yypt-0].str
		}
	case 102:
		//line cc.y:761
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.str = yyS[yypt-0].str
		}
	case 103:
		//line cc.y:766
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.str = yyS[yypt-0].str
		}
	case 104:
		//line cc.y:771
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.str = yyS[yypt-0].str
		}
	case 105:
		//line cc.y:776
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.str = yyS[yypt-0].str
		}
	case 106:
		//line cc.y:784
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.str = yyS[yypt-0].str
		}
	case 107:
		//line cc.y:789
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.str = yyS[yypt-0].str
		}
	case 108:
		//line cc.y:797
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.str = yyS[yypt-0].str
		}
	case 109:
		//line cc.y:802
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.str = yyS[yypt-0].str
		}
	case 110:
		//line cc.y:807
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.str = yyS[yypt-0].str
		}
	case 111:
		//line cc.y:812
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.str = yyS[yypt-0].str
		}
	case 112:
		//line cc.y:817
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.str = yyS[yypt-0].str
		}
	case 113:
		//line cc.y:822
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.str = yyS[yypt-0].str
		}
	case 114:
		//line cc.y:827
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.str = yyS[yypt-0].str
		}
	case 115:
		//line cc.y:832
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.str = yyS[yypt-0].str
		}
	case 116:
		//line cc.y:837
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.str = yyS[yypt-0].str
		}
	case 117:
		//line cc.y:844
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.str = yyS[yypt-0].str
		}
	case 118:
		//line cc.y:849
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.str = yyS[yypt-0].str
		}
	case 119:
		//line cc.y:856
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.str = yyS[yypt-0].str
		}
	case 120:
		//line cc.y:861
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.str = yyS[yypt-0].str
		}
	case 121:
		//line cc.y:869
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.typ = yyS[yypt-0].typ
			if yyVAL.typ == nil {
				yyVAL.typ = &Type{Kind: TypedefType, Name: yyS[yypt-0].str}
			}
		}
	case 122:
		//line cc.y:885
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.tc.c, yyVAL.tc.q, yyVAL.tc.t = splitTypeWords(append(yyS[yypt-0].strs, "int"))
		}
	case 123:
		//line cc.y:890
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.tc.c, yyVAL.tc.q, _ = splitTypeWords(append(yyS[yypt-2].strs, yyS[yypt-0].strs...))
			yyVAL.tc.t = yyS[yypt-1].typ
		}
	case 124:
		//line cc.y:896
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyS[yypt-2].strs = append(yyS[yypt-2].strs, yyS[yypt-1].str)
			yyS[yypt-2].strs = append(yyS[yypt-2].strs, yyS[yypt-0].strs...)
			yyVAL.tc.c, yyVAL.tc.q, yyVAL.tc.t = splitTypeWords(yyS[yypt-2].strs)
		}
	case 125:
		//line cc.y:903
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.tc.c, yyVAL.tc.q, _ = splitTypeWords(yyS[yypt-0].strs)
			yyVAL.tc.t = yyS[yypt-1].typ
		}
	case 126:
		//line cc.y:909
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			var ts []string
			ts = append(ts, yyS[yypt-1].str)
			ts = append(ts, yyS[yypt-0].strs...)
			yyVAL.tc.c, yyVAL.tc.q, yyVAL.tc.t = splitTypeWords(ts)
		}
	case 127:
		//line cc.y:920
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
	case 128:
		//line cc.y:933
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.typ = yyS[yypt-0].abdecor(yyS[yypt-1].typ)
		}
	case 129:
		//line cc.y:941
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			// TODO: use $1.q
			yyVAL.span = Span{}
			yyVAL.decls = nil
			for _, idec := range yyS[yypt-1].idecs {
				typ, name := idec.d(yyS[yypt-2].tc.t)
				decl := &Decl{Span: yyVAL.span, Name: name, Type: typ, Storage: yyS[yypt-2].tc.c, Init: idec.i}
				yylex.(*lexer).pushDecl(decl)
				yyVAL.decls = append(yyVAL.decls, decl)
			}
			if yyS[yypt-1].idecs == nil {
				decl := &Decl{Span: yyVAL.span, Name: "", Type: yyS[yypt-2].tc.t, Storage: yyS[yypt-2].tc.c}
				yylex.(*lexer).pushDecl(decl)
				yyVAL.decls = append(yyVAL.decls, decl)
			}
		}
	case 130:
		//line cc.y:962
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.decls = yyS[yypt-0].decls
		}
	case 131:
		//line cc.y:967
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.decls = []*Decl{yyS[yypt-0].decl}
		}
	case 132:
		//line cc.y:972
		{
			yyVAL.decls = yyS[yypt-1].decls
		}
	case 133:
		//line cc.y:978
		{
			yyVAL.span = span(yyS[yypt-3].span, yyS[yypt-0].span)
			typ, name := yyS[yypt-2].decor(yyS[yypt-3].tc.t)
			yyVAL.decl = &Decl{Span: yyVAL.span, Name: name, Type: typ}
			if yyS[yypt-1].decls != nil {
				yylex.(*lexer).Errorf("cannot use pre-prototype definitions")
			}
			yyVAL.decl.Body = yyS[yypt-0].stmt
			yylex.(*lexer).pushDecl(yyVAL.decl)
		}
	case 134:
		//line cc.y:991
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.str = yyS[yypt-0].str
		}
	case 135:
		//line cc.y:996
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.str = yyS[yypt-0].str
		}
	case 136:
		//line cc.y:1004
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.tk = Struct
		}
	case 137:
		//line cc.y:1009
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.tk = Union
		}
	case 138:
		//line cc.y:1016
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.decor = yyS[yypt-0].decor
		}
	case 139:
		//line cc.y:1021
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			name := yyS[yypt-2].str
			expr := yyS[yypt-0].expr
			yyVAL.decor = func(t *Type) (*Type, string) {
				t.Width = expr
				return t, name
			}
		}
	case 140:
		//line cc.y:1033
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
	case 141:
		//line cc.y:1047
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.typ = yylex.(*lexer).pushType(&Type{Span: yyVAL.span, Kind: yyS[yypt-1].tk, Tag: yyS[yypt-0].str})
		}
	case 142:
		//line cc.y:1052
		{
			yyVAL.span = span(yyS[yypt-4].span, yyS[yypt-0].span)
			yyVAL.typ = yylex.(*lexer).pushType(&Type{Span: yyVAL.span, Kind: yyS[yypt-4].tk, Tag: yyS[yypt-3].str, Decls: yyS[yypt-1].decls})
		}
	case 143:
		//line cc.y:1059
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.prefix = &Prefix{Span: yyVAL.span, Dot: yyS[yypt-0].str}
		}
	case 144:
		//line cc.y:1066
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: Arrow, Left: yyS[yypt-2].expr, Text: yyS[yypt-0].str}
		}
	case 145:
		//line cc.y:1071
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.expr = &Expr{Span: yyVAL.span, Op: Dot, Left: yyS[yypt-2].expr, Text: yyS[yypt-0].str}
		}
	case 146:
		//line cc.y:1079
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.typ = yylex.(*lexer).pushType(&Type{Span: yyVAL.span, Kind: Enum, Tag: yyS[yypt-0].str})
		}
	case 147:
		//line cc.y:1084
		{
			yyVAL.span = span(yyS[yypt-5].span, yyS[yypt-0].span)
			yyVAL.typ = yylex.(*lexer).pushType(&Type{Span: yyVAL.span, Kind: Enum, Tag: yyS[yypt-4].str, Decls: yyS[yypt-2].decls})
		}
	case 148:
		//line cc.y:1091
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			var x *Init
			if yyS[yypt-0].expr != nil {
				x = &Init{Span: yyVAL.span, Expr: yyS[yypt-0].expr}
			}
			yyVAL.decl = &Decl{Span: yyVAL.span, Name: yyS[yypt-1].str, Init: x}
		}
	case 149:
		//line cc.y:1102
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.expr = yyS[yypt-0].expr
		}
	case 150:
		//line cc.y:1110
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.init = &Init{Span: yyVAL.span, Expr: yyS[yypt-0].expr}
		}
	case 151:
		//line cc.y:1115
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.init = &Init{Span: yyVAL.span, Braced: yyS[yypt-0].inits}
		}
	case 152:
		//line cc.y:1122
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.inits = []*Init{}
		}
	case 153:
		//line cc.y:1127
		{
			yyVAL.span = span(yyS[yypt-3].span, yyS[yypt-0].span)
			yyVAL.inits = append(yyS[yypt-2].inits, yyS[yypt-1].init)
		}
	case 154:
		//line cc.y:1132
		{
			yyVAL.span = span(yyS[yypt-4].span, yyS[yypt-0].span)
			yyVAL.inits = append(yyS[yypt-3].inits, yyS[yypt-2].init)
		}
	case 155:
		//line cc.y:1138
		{
			yyVAL.span = Span{}
			yyVAL.inits = nil
		}
	case 156:
		//line cc.y:1143
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.inits = append(yyS[yypt-2].inits, yyS[yypt-1].init)
		}
	case 157:
		//line cc.y:1150
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.init = yyS[yypt-0].init
		}
	case 158:
		//line cc.y:1155
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.init = yyS[yypt-0].init
			yyVAL.init.Prefix = yyS[yypt-2].prefixes
		}
	case 159:
		//line cc.y:1163
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.prefix = &Prefix{Span: yyVAL.span, Index: yyS[yypt-1].expr}
		}
	case 160:
		//line cc.y:1169
		{
			yyVAL.span = Span{}
		}
	case 161:
		//line cc.y:1173
		{
			yyVAL.span = yyS[yypt-0].span
		}
	case 162:
		//line cc.y:1178
		{
			yyVAL.span = Span{}
		}
	case 163:
		//line cc.y:1182
		{
			yyVAL.span = yyS[yypt-0].span
		}
	case 164:
		//line cc.y:1191
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.prefixes = []*Prefix{yyS[yypt-0].prefix}
		}
	case 165:
		//line cc.y:1196
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.prefixes = append(yyS[yypt-1].prefixes, yyS[yypt-0].prefix)
		}
	case 166:
		//line cc.y:1202
		{
			yyVAL.span = Span{}
			yyVAL.str = ""
		}
	case 167:
		//line cc.y:1207
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.str = yyS[yypt-0].str
		}
	case 168:
		//line cc.y:1213
		{
			yyVAL.span = Span{}
			yyVAL.expr = nil
		}
	case 169:
		//line cc.y:1218
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.expr = yyS[yypt-0].expr
		}
	case 170:
		//line cc.y:1224
		{
			yyVAL.span = Span{}
			yyVAL.expr = nil
		}
	case 171:
		//line cc.y:1229
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.expr = yyS[yypt-0].expr
		}
	case 172:
		//line cc.y:1236
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.exprs = []*Expr{yyS[yypt-0].expr}
		}
	case 173:
		//line cc.y:1241
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.exprs = append(yyS[yypt-2].exprs, yyS[yypt-0].expr)
		}
	case 174:
		//line cc.y:1247
		{
			yyVAL.span = Span{}
			yyVAL.exprs = nil
		}
	case 175:
		//line cc.y:1252
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.exprs = yyS[yypt-0].exprs
		}
	case 176:
		//line cc.y:1258
		{
			yyVAL.span = Span{}
			yyVAL.decls = nil
		}
	case 177:
		//line cc.y:1263
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.decls = append(yyS[yypt-1].decls, yyS[yypt-0].decls...)
		}
	case 178:
		//line cc.y:1269
		{
			yyVAL.span = Span{}
			yyVAL.labels = nil
		}
	case 179:
		//line cc.y:1274
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.labels = append(yyS[yypt-1].labels, yyS[yypt-0].label)
		}
	case 180:
		//line cc.y:1281
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.decls = []*Decl{yyS[yypt-0].decl}
		}
	case 181:
		//line cc.y:1286
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.decls = append(yyS[yypt-2].decls, yyS[yypt-0].decl)
		}
	case 182:
		//line cc.y:1292
		{
			yyVAL.span = Span{}
			yyVAL.decls = nil
		}
	case 183:
		//line cc.y:1297
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.decls = yyS[yypt-0].decls
		}
	case 184:
		//line cc.y:1304
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.idecs = []idecor{yyS[yypt-0].idec}
		}
	case 185:
		//line cc.y:1309
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.idecs = append(yyS[yypt-2].idecs, yyS[yypt-0].idec)
		}
	case 186:
		//line cc.y:1315
		{
			yyVAL.span = Span{}
			yyVAL.idecs = nil
		}
	case 187:
		//line cc.y:1320
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.idecs = yyS[yypt-0].idecs
		}
	case 188:
		//line cc.y:1327
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.strs = []string{yyS[yypt-0].str}
		}
	case 189:
		//line cc.y:1332
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.strs = append(yyS[yypt-1].strs, yyS[yypt-0].str)
		}
	case 190:
		//line cc.y:1338
		{
			yyVAL.span = Span{}
			yyVAL.strs = nil
		}
	case 191:
		//line cc.y:1343
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.strs = yyS[yypt-0].strs
		}
	case 192:
		//line cc.y:1350
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.strs = []string{yyS[yypt-0].str}
		}
	case 193:
		//line cc.y:1355
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.strs = append(yyS[yypt-1].strs, yyS[yypt-0].str)
		}
	case 194:
		//line cc.y:1361
		{
			yyVAL.span = Span{}
			yyVAL.strs = nil
		}
	case 195:
		//line cc.y:1366
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.strs = yyS[yypt-0].strs
		}
	case 196:
		//line cc.y:1373
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.strs = []string{yyS[yypt-0].str}
		}
	case 197:
		//line cc.y:1378
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.strs = append(yyS[yypt-1].strs, yyS[yypt-0].str)
		}
	case 198:
		//line cc.y:1384
		{
			yyVAL.span = Span{}
			yyVAL.strs = nil
		}
	case 199:
		//line cc.y:1389
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.strs = yyS[yypt-0].strs
		}
	case 200:
		//line cc.y:1396
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.decors = nil
			yyVAL.decors = append(yyVAL.decors, yyS[yypt-0].decor)
		}
	case 201:
		//line cc.y:1402
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.decors = append(yyS[yypt-2].decors, yyS[yypt-0].decor)
		}
	case 202:
		//line cc.y:1408
		{
			yyVAL.span = Span{}
			yyVAL.decors = nil
		}
	case 203:
		//line cc.y:1413
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.decors = yyS[yypt-0].decors
		}
	case 204:
		//line cc.y:1420
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.decls = yyS[yypt-0].decls
		}
	case 205:
		//line cc.y:1425
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.decls = append(yyS[yypt-1].decls, yyS[yypt-0].decls...)
		}
	case 206:
		//line cc.y:1431
		{
			yyVAL.span = Span{}
			yyVAL.expr = nil
		}
	case 207:
		//line cc.y:1436
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.expr = yyS[yypt-0].expr
		}
	case 208:
		//line cc.y:1443
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.decls = []*Decl{yyS[yypt-0].decl}
		}
	case 209:
		//line cc.y:1448
		{
			yyVAL.span = span(yyS[yypt-2].span, yyS[yypt-0].span)
			yyVAL.decls = append(yyS[yypt-2].decls, yyS[yypt-0].decl)
		}
	case 210:
		//line cc.y:1455
		{
			yyVAL.span = yyS[yypt-0].span
			yyVAL.strs = []string{yyS[yypt-0].str}
		}
	case 211:
		//line cc.y:1460
		{
			yyVAL.span = span(yyS[yypt-1].span, yyS[yypt-0].span)
			yyVAL.strs = append(yyS[yypt-1].strs, yyS[yypt-0].str)
		}
	}
	goto yystack /* stack new state and value */
}
