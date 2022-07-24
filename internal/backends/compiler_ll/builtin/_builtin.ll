; ModuleID = '_builtin.c'
source_filename = "_builtin.c"
target datalayout = "e-m:o-p270:32:32-p271:32:32-p272:64:64-i64:64-f80:128-n8:16:32:64-S128"
target triple = "x86_64-apple-macosx12.0.0"

%struct.ugo_string_t = type { i8*, i32 }
%struct.__va_list_tag = type { i32, i32, i8*, i8* }

@.str = private unnamed_addr constant [3 x i8] c"%c\00", align 1
@.str.1 = private unnamed_addr constant [3 x i8] c"%s\00", align 1
@.str.2 = private unnamed_addr constant [5 x i8] c"%.*s\00", align 1
@.str.3 = private unnamed_addr constant [5 x i8] c"true\00", align 1
@.str.4 = private unnamed_addr constant [6 x i8] c"false\00", align 1
@.str.5 = private unnamed_addr constant [3 x i8] c"%d\00", align 1
@.str.6 = private unnamed_addr constant [5 x i8] c"%lld\00", align 1
@.str.7 = private unnamed_addr constant [5 x i8] c"0x%x\00", align 1
@.str.8 = private unnamed_addr constant [4 x i8] c"%d\0A\00", align 1

; Function Attrs: noinline nounwind optnone ssp uwtable
define %struct.ugo_string_t* @ugo_string_new(i32 %0, i8* %1) #0 {
  %3 = alloca i32, align 4
  %4 = alloca i8*, align 8
  %5 = alloca %struct.ugo_string_t*, align 8
  store i32 %0, i32* %3, align 4
  store i8* %1, i8** %4, align 8
  %6 = load i32, i32* %3, align 4
  %7 = sext i32 %6 to i64
  %8 = add i64 16, %7
  %9 = add i64 %8, 1
  %10 = call i8* @malloc(i64 %9) #7
  %11 = bitcast i8* %10 to %struct.ugo_string_t*
  store %struct.ugo_string_t* %11, %struct.ugo_string_t** %5, align 8
  %12 = load %struct.ugo_string_t*, %struct.ugo_string_t** %5, align 8
  %13 = getelementptr inbounds %struct.ugo_string_t, %struct.ugo_string_t* %12, i64 1
  %14 = bitcast %struct.ugo_string_t* %13 to i8*
  %15 = load %struct.ugo_string_t*, %struct.ugo_string_t** %5, align 8
  %16 = getelementptr inbounds %struct.ugo_string_t, %struct.ugo_string_t* %15, i32 0, i32 0
  store i8* %14, i8** %16, align 8
  %17 = load i32, i32* %3, align 4
  %18 = load %struct.ugo_string_t*, %struct.ugo_string_t** %5, align 8
  %19 = getelementptr inbounds %struct.ugo_string_t, %struct.ugo_string_t* %18, i32 0, i32 1
  store i32 %17, i32* %19, align 8
  %20 = load i8*, i8** %4, align 8
  %21 = icmp ne i8* %20, null
  br i1 %21, label %22, label %40

22:                                               ; preds = %2
  %23 = load %struct.ugo_string_t*, %struct.ugo_string_t** %5, align 8
  %24 = getelementptr inbounds %struct.ugo_string_t, %struct.ugo_string_t* %23, i32 0, i32 0
  %25 = load i8*, i8** %24, align 8
  %26 = load i8*, i8** %4, align 8
  %27 = load i32, i32* %3, align 4
  %28 = sext i32 %27 to i64
  %29 = load %struct.ugo_string_t*, %struct.ugo_string_t** %5, align 8
  %30 = getelementptr inbounds %struct.ugo_string_t, %struct.ugo_string_t* %29, i32 0, i32 0
  %31 = load i8*, i8** %30, align 8
  %32 = call i64 @llvm.objectsize.i64.p0i8(i8* %31, i1 false, i1 true, i1 false)
  %33 = call i8* @__memcpy_chk(i8* %25, i8* %26, i64 %28, i64 %32) #8
  %34 = load %struct.ugo_string_t*, %struct.ugo_string_t** %5, align 8
  %35 = getelementptr inbounds %struct.ugo_string_t, %struct.ugo_string_t* %34, i32 0, i32 0
  %36 = load i8*, i8** %35, align 8
  %37 = load i32, i32* %3, align 4
  %38 = sext i32 %37 to i64
  %39 = getelementptr inbounds i8, i8* %36, i64 %38
  store i8 0, i8* %39, align 1
  br label %52

40:                                               ; preds = %2
  %41 = load %struct.ugo_string_t*, %struct.ugo_string_t** %5, align 8
  %42 = getelementptr inbounds %struct.ugo_string_t, %struct.ugo_string_t* %41, i32 0, i32 0
  %43 = load i8*, i8** %42, align 8
  %44 = load i32, i32* %3, align 4
  %45 = add nsw i32 %44, 1
  %46 = sext i32 %45 to i64
  %47 = load %struct.ugo_string_t*, %struct.ugo_string_t** %5, align 8
  %48 = getelementptr inbounds %struct.ugo_string_t, %struct.ugo_string_t* %47, i32 0, i32 0
  %49 = load i8*, i8** %48, align 8
  %50 = call i64 @llvm.objectsize.i64.p0i8(i8* %49, i1 false, i1 true, i1 false)
  %51 = call i8* @__memset_chk(i8* %43, i32 0, i64 %46, i64 %50) #8
  br label %52

52:                                               ; preds = %40, %22
  %53 = load %struct.ugo_string_t*, %struct.ugo_string_t** %5, align 8
  ret %struct.ugo_string_t* %53
}

; Function Attrs: allocsize(0)
declare i8* @malloc(i64) #1

; Function Attrs: nounwind
declare i8* @__memcpy_chk(i8*, i8*, i64, i64) #2

; Function Attrs: nofree nosync nounwind readnone speculatable willreturn
declare i64 @llvm.objectsize.i64.p0i8(i8*, i1 immarg, i1 immarg, i1 immarg) #3

; Function Attrs: nounwind
declare i8* @__memset_chk(i8*, i32, i64, i64) #2

; Function Attrs: noinline nounwind optnone ssp uwtable
define %struct.ugo_string_t* @ugo_string_clone(%struct.ugo_string_t* %0) #0 {
  %2 = alloca %struct.ugo_string_t*, align 8
  %3 = alloca %struct.ugo_string_t*, align 8
  store %struct.ugo_string_t* %0, %struct.ugo_string_t** %2, align 8
  %4 = load %struct.ugo_string_t*, %struct.ugo_string_t** %2, align 8
  %5 = getelementptr inbounds %struct.ugo_string_t, %struct.ugo_string_t* %4, i32 0, i32 1
  %6 = load i32, i32* %5, align 8
  %7 = sext i32 %6 to i64
  %8 = add i64 16, %7
  %9 = add i64 %8, 1
  %10 = call i8* @malloc(i64 %9) #7
  %11 = bitcast i8* %10 to %struct.ugo_string_t*
  store %struct.ugo_string_t* %11, %struct.ugo_string_t** %3, align 8
  %12 = load %struct.ugo_string_t*, %struct.ugo_string_t** %3, align 8
  %13 = getelementptr inbounds %struct.ugo_string_t, %struct.ugo_string_t* %12, i32 0, i32 0
  %14 = load i8*, i8** %13, align 8
  %15 = load %struct.ugo_string_t*, %struct.ugo_string_t** %2, align 8
  %16 = getelementptr inbounds %struct.ugo_string_t, %struct.ugo_string_t* %15, i32 0, i32 0
  %17 = load i8*, i8** %16, align 8
  %18 = load %struct.ugo_string_t*, %struct.ugo_string_t** %2, align 8
  %19 = getelementptr inbounds %struct.ugo_string_t, %struct.ugo_string_t* %18, i32 0, i32 1
  %20 = load i32, i32* %19, align 8
  %21 = sext i32 %20 to i64
  %22 = load %struct.ugo_string_t*, %struct.ugo_string_t** %3, align 8
  %23 = getelementptr inbounds %struct.ugo_string_t, %struct.ugo_string_t* %22, i32 0, i32 0
  %24 = load i8*, i8** %23, align 8
  %25 = call i64 @llvm.objectsize.i64.p0i8(i8* %24, i1 false, i1 true, i1 false)
  %26 = call i8* @__memcpy_chk(i8* %14, i8* %17, i64 %21, i64 %25) #8
  %27 = load %struct.ugo_string_t*, %struct.ugo_string_t** %3, align 8
  %28 = getelementptr inbounds %struct.ugo_string_t, %struct.ugo_string_t* %27, i32 0, i32 0
  %29 = load i8*, i8** %28, align 8
  %30 = load %struct.ugo_string_t*, %struct.ugo_string_t** %2, align 8
  %31 = getelementptr inbounds %struct.ugo_string_t, %struct.ugo_string_t* %30, i32 0, i32 1
  %32 = load i32, i32* %31, align 8
  %33 = sext i32 %32 to i64
  %34 = getelementptr inbounds i8, i8* %29, i64 %33
  store i8 0, i8* %34, align 1
  %35 = load %struct.ugo_string_t*, %struct.ugo_string_t** %2, align 8
  %36 = getelementptr inbounds %struct.ugo_string_t, %struct.ugo_string_t* %35, i32 0, i32 1
  %37 = load i32, i32* %36, align 8
  %38 = load %struct.ugo_string_t*, %struct.ugo_string_t** %3, align 8
  %39 = getelementptr inbounds %struct.ugo_string_t, %struct.ugo_string_t* %38, i32 0, i32 1
  store i32 %37, i32* %39, align 8
  %40 = load %struct.ugo_string_t*, %struct.ugo_string_t** %3, align 8
  ret %struct.ugo_string_t* %40
}

; Function Attrs: noinline nounwind optnone ssp uwtable
define void @ugo_string_free(%struct.ugo_string_t* %0) #0 {
  %2 = alloca %struct.ugo_string_t*, align 8
  store %struct.ugo_string_t* %0, %struct.ugo_string_t** %2, align 8
  %3 = load %struct.ugo_string_t*, %struct.ugo_string_t** %2, align 8
  %4 = bitcast %struct.ugo_string_t* %3 to i8*
  call void @free(i8* %4)
  ret void
}

declare void @free(i8*) #4

; Function Attrs: noinline nounwind optnone ssp uwtable
define i8* @ugo_cstring_join(i8* %0, i8* %1) #0 {
  %3 = alloca i8*, align 8
  %4 = alloca i8*, align 8
  %5 = alloca i32, align 4
  %6 = alloca i32, align 4
  %7 = alloca i8*, align 8
  store i8* %0, i8** %3, align 8
  store i8* %1, i8** %4, align 8
  %8 = load i8*, i8** %3, align 8
  %9 = call i64 @strlen(i8* %8)
  %10 = trunc i64 %9 to i32
  store i32 %10, i32* %5, align 4
  %11 = load i8*, i8** %4, align 8
  %12 = call i64 @strlen(i8* %11)
  %13 = trunc i64 %12 to i32
  store i32 %13, i32* %6, align 4
  %14 = load i32, i32* %5, align 4
  %15 = load i32, i32* %6, align 4
  %16 = add nsw i32 %14, %15
  %17 = add nsw i32 %16, 1
  %18 = sext i32 %17 to i64
  %19 = call i8* @malloc(i64 %18) #7
  store i8* %19, i8** %7, align 8
  %20 = load i8*, i8** %7, align 8
  %21 = load i8*, i8** %3, align 8
  %22 = load i8*, i8** %7, align 8
  %23 = call i64 @llvm.objectsize.i64.p0i8(i8* %22, i1 false, i1 true, i1 false)
  %24 = call i8* @__strcpy_chk(i8* %20, i8* %21, i64 %23) #8
  %25 = load i8*, i8** %7, align 8
  %26 = load i8*, i8** %4, align 8
  %27 = load i8*, i8** %7, align 8
  %28 = call i64 @llvm.objectsize.i64.p0i8(i8* %27, i1 false, i1 true, i1 false)
  %29 = call i8* @__strcat_chk(i8* %25, i8* %26, i64 %28) #8
  %30 = load i8*, i8** %7, align 8
  ret i8* %30
}

declare i64 @strlen(i8*) #4

; Function Attrs: nounwind
declare i8* @__strcpy_chk(i8*, i8*, i64) #2

; Function Attrs: nounwind
declare i8* @__strcat_chk(i8*, i8*, i64) #2

; Function Attrs: noinline nounwind optnone ssp uwtable
define i8* @ugo_cstring_slice(i8* %0, i32 %1, i32 %2) #0 {
  %4 = alloca i8*, align 8
  %5 = alloca i32, align 4
  %6 = alloca i32, align 4
  %7 = alloca i32, align 4
  %8 = alloca i32, align 4
  %9 = alloca i8*, align 8
  store i8* %0, i8** %4, align 8
  store i32 %1, i32* %5, align 4
  store i32 %2, i32* %6, align 4
  %10 = load i8*, i8** %4, align 8
  %11 = call i64 @strlen(i8* %10)
  %12 = trunc i64 %11 to i32
  store i32 %12, i32* %7, align 4
  store i32 0, i32* %8, align 4
  store i8* null, i8** %9, align 8
  %13 = load i32, i32* %5, align 4
  %14 = icmp slt i32 %13, 0
  br i1 %14, label %15, label %16

15:                                               ; preds = %3
  store i32 0, i32* %5, align 4
  br label %16

16:                                               ; preds = %15, %3
  %17 = load i32, i32* %6, align 4
  %18 = icmp slt i32 %17, 0
  br i1 %18, label %19, label %21

19:                                               ; preds = %16
  %20 = load i32, i32* %7, align 4
  store i32 %20, i32* %6, align 4
  br label %21

21:                                               ; preds = %19, %16
  %22 = load i32, i32* %5, align 4
  %23 = load i32, i32* %7, align 4
  %24 = icmp sge i32 %22, %23
  br i1 %24, label %25, label %28

25:                                               ; preds = %21
  %26 = load i32, i32* %7, align 4
  %27 = sub nsw i32 %26, 1
  store i32 %27, i32* %5, align 4
  br label %28

28:                                               ; preds = %25, %21
  %29 = load i32, i32* %6, align 4
  %30 = load i32, i32* %7, align 4
  %31 = icmp sgt i32 %29, %30
  br i1 %31, label %32, label %34

32:                                               ; preds = %28
  %33 = load i32, i32* %7, align 4
  store i32 %33, i32* %5, align 4
  br label %34

34:                                               ; preds = %32, %28
  %35 = load i32, i32* %6, align 4
  %36 = load i32, i32* %5, align 4
  %37 = sub nsw i32 %35, %36
  store i32 %37, i32* %8, align 4
  %38 = load i32, i32* %8, align 4
  %39 = add nsw i32 %38, 1
  %40 = sext i32 %39 to i64
  %41 = call i8* @malloc(i64 %40) #7
  store i8* %41, i8** %9, align 8
  %42 = load i8*, i8** %9, align 8
  %43 = load i8*, i8** %4, align 8
  %44 = load i32, i32* %5, align 4
  %45 = sext i32 %44 to i64
  %46 = getelementptr inbounds i8, i8* %43, i64 %45
  %47 = load i32, i32* %8, align 4
  %48 = sext i32 %47 to i64
  %49 = load i8*, i8** %9, align 8
  %50 = call i64 @llvm.objectsize.i64.p0i8(i8* %49, i1 false, i1 true, i1 false)
  %51 = call i8* @__memcpy_chk(i8* %42, i8* %46, i64 %48, i64 %50) #8
  %52 = load i8*, i8** %9, align 8
  %53 = load i32, i32* %8, align 4
  %54 = sext i32 %53 to i64
  %55 = getelementptr inbounds i8, i8* %52, i64 %54
  store i8 0, i8* %55, align 1
  %56 = load i8*, i8** %9, align 8
  ret i8* %56
}

; Function Attrs: noinline nounwind optnone ssp uwtable
define i32 @ugo_cstring_index(i8* %0, i32 %1) #0 {
  %3 = alloca i8*, align 8
  %4 = alloca i32, align 4
  store i8* %0, i8** %3, align 8
  store i32 %1, i32* %4, align 4
  %5 = load i8*, i8** %3, align 8
  %6 = load i32, i32* %4, align 4
  %7 = sext i32 %6 to i64
  %8 = getelementptr inbounds i8, i8* %5, i64 %7
  %9 = load i8, i8* %8, align 1
  %10 = sext i8 %9 to i32
  ret i32 %10
}

; Function Attrs: noinline nounwind optnone ssp uwtable
define i32 @ugo_cstring_cmp(i8* %0, i8* %1) #0 {
  %3 = alloca i8*, align 8
  %4 = alloca i8*, align 8
  store i8* %0, i8** %3, align 8
  store i8* %1, i8** %4, align 8
  %5 = load i8*, i8** %3, align 8
  %6 = load i8*, i8** %4, align 8
  %7 = call i32 @strcmp(i8* %5, i8* %6)
  ret i32 %7
}

declare i32 @strcmp(i8*, i8*) #4

; Function Attrs: noinline nounwind optnone ssp uwtable
define i32 @ugo_print_rune(i32 %0) #0 {
  %2 = alloca i32, align 4
  store i32 %0, i32* %2, align 4
  %3 = load i32, i32* %2, align 4
  %4 = call i32 (i8*, ...) @printf(i8* getelementptr inbounds ([3 x i8], [3 x i8]* @.str, i64 0, i64 0), i32 %3)
  ret i32 %4
}

declare i32 @printf(i8*, ...) #4

; Function Attrs: noinline nounwind optnone ssp uwtable
define i32 @ugo_print_cstring(i8* %0) #0 {
  %2 = alloca i8*, align 8
  store i8* %0, i8** %2, align 8
  %3 = load i8*, i8** %2, align 8
  %4 = call i32 (i8*, ...) @printf(i8* getelementptr inbounds ([3 x i8], [3 x i8]* @.str.1, i64 0, i64 0), i8* %3)
  ret i32 %4
}

; Function Attrs: noinline nounwind optnone ssp uwtable
define i32 @ugo_print_cstring_len(i8* %0, i32 %1) #0 {
  %3 = alloca i8*, align 8
  %4 = alloca i32, align 4
  store i8* %0, i8** %3, align 8
  store i32 %1, i32* %4, align 4
  br label %5

5:                                                ; preds = %9, %2
  %6 = load i32, i32* %4, align 4
  %7 = add nsw i32 %6, -1
  store i32 %7, i32* %4, align 4
  %8 = icmp sgt i32 %6, 0
  br i1 %8, label %9, label %15

9:                                                ; preds = %5
  %10 = load i8*, i8** %3, align 8
  %11 = getelementptr inbounds i8, i8* %10, i32 1
  store i8* %11, i8** %3, align 8
  %12 = load i8, i8* %10, align 1
  %13 = sext i8 %12 to i32
  %14 = call i32 (i8*, ...) @printf(i8* getelementptr inbounds ([3 x i8], [3 x i8]* @.str, i64 0, i64 0), i32 %13)
  br label %5, !llvm.loop !4

15:                                               ; preds = %5
  ret i32 1
}

; Function Attrs: noinline nounwind optnone ssp uwtable
define i32 @ugo_print_string(%struct.ugo_string_t* %0) #0 {
  %2 = alloca %struct.ugo_string_t*, align 8
  store %struct.ugo_string_t* %0, %struct.ugo_string_t** %2, align 8
  %3 = load %struct.ugo_string_t*, %struct.ugo_string_t** %2, align 8
  %4 = getelementptr inbounds %struct.ugo_string_t, %struct.ugo_string_t* %3, i32 0, i32 1
  %5 = load i32, i32* %4, align 8
  %6 = load %struct.ugo_string_t*, %struct.ugo_string_t** %2, align 8
  %7 = getelementptr inbounds %struct.ugo_string_t, %struct.ugo_string_t* %6, i32 0, i32 0
  %8 = load i8*, i8** %7, align 8
  %9 = call i32 (i8*, ...) @printf(i8* getelementptr inbounds ([5 x i8], [5 x i8]* @.str.2, i64 0, i64 0), i32 %5, i8* %8)
  ret i32 %9
}

; Function Attrs: noinline nounwind optnone ssp uwtable
define i32 @ugo_print_bool(i8 zeroext %0) #0 {
  %2 = alloca i8, align 1
  store i8 %0, i8* %2, align 1
  %3 = load i8, i8* %2, align 1
  %4 = zext i8 %3 to i32
  %5 = icmp ne i32 %4, 0
  %6 = zext i1 %5 to i64
  %7 = select i1 %5, i8* getelementptr inbounds ([5 x i8], [5 x i8]* @.str.3, i64 0, i64 0), i8* getelementptr inbounds ([6 x i8], [6 x i8]* @.str.4, i64 0, i64 0)
  %8 = call i32 (i8*, ...) @printf(i8* getelementptr inbounds ([3 x i8], [3 x i8]* @.str.1, i64 0, i64 0), i8* %7)
  ret i32 %8
}

; Function Attrs: noinline nounwind optnone ssp uwtable
define i32 @ugo_print_int(i32 %0) #0 {
  %2 = alloca i32, align 4
  store i32 %0, i32* %2, align 4
  %3 = load i32, i32* %2, align 4
  %4 = call i32 (i8*, ...) @printf(i8* getelementptr inbounds ([3 x i8], [3 x i8]* @.str.5, i64 0, i64 0), i32 %3)
  ret i32 %4
}

; Function Attrs: noinline nounwind optnone ssp uwtable
define i32 @ugo_print_int64(i64 %0) #0 {
  %2 = alloca i64, align 8
  store i64 %0, i64* %2, align 8
  %3 = load i64, i64* %2, align 8
  %4 = call i32 (i8*, ...) @printf(i8* getelementptr inbounds ([5 x i8], [5 x i8]* @.str.6, i64 0, i64 0), i64 %3)
  ret i32 %4
}

; Function Attrs: noinline nounwind optnone ssp uwtable
define i32 @ugo_print_ptr(i8* %0) #0 {
  %2 = alloca i8*, align 8
  store i8* %0, i8** %2, align 8
  %3 = load i8*, i8** %2, align 8
  %4 = load i8*, i8** %2, align 8
  %5 = call i32 (i8*, ...) @printf(i8* getelementptr inbounds ([5 x i8], [5 x i8]* @.str.7, i64 0, i64 0), i8* %3, i8* %4)
  ret i32 %5
}

; Function Attrs: noinline nounwind optnone ssp uwtable
define i32 @ugo_printf(i8* %0, ...) #0 {
  %2 = alloca i8*, align 8
  %3 = alloca [1 x %struct.__va_list_tag], align 16
  %4 = alloca i32, align 4
  store i8* %0, i8** %2, align 8
  %5 = getelementptr inbounds [1 x %struct.__va_list_tag], [1 x %struct.__va_list_tag]* %3, i64 0, i64 0
  %6 = bitcast %struct.__va_list_tag* %5 to i8*
  call void @llvm.va_start(i8* %6)
  %7 = load i8*, i8** %2, align 8
  %8 = getelementptr inbounds [1 x %struct.__va_list_tag], [1 x %struct.__va_list_tag]* %3, i64 0, i64 0
  %9 = call i32 @vprintf(i8* %7, %struct.__va_list_tag* %8)
  store i32 %9, i32* %4, align 4
  %10 = getelementptr inbounds [1 x %struct.__va_list_tag], [1 x %struct.__va_list_tag]* %3, i64 0, i64 0
  %11 = bitcast %struct.__va_list_tag* %10 to i8*
  call void @llvm.va_end(i8* %11)
  %12 = load i32, i32* %4, align 4
  ret i32 %12
}

; Function Attrs: nofree nosync nounwind willreturn
declare void @llvm.va_start(i8*) #5

declare i32 @vprintf(i8*, %struct.__va_list_tag*) #4

; Function Attrs: nofree nosync nounwind willreturn
declare void @llvm.va_end(i8*) #5

; Function Attrs: noinline nounwind optnone ssp uwtable
define i32 @ugo_builtin_println(i32 %0) #0 {
  %2 = alloca i32, align 4
  store i32 %0, i32* %2, align 4
  %3 = load i32, i32* %2, align 4
  %4 = call i32 (i8*, ...) @printf(i8* getelementptr inbounds ([4 x i8], [4 x i8]* @.str.8, i64 0, i64 0), i32 %3)
  ret i32 %4
}

; Function Attrs: noinline nounwind optnone ssp uwtable
define i32 @ugo_builtin_exit(i32 %0) #0 {
  %2 = alloca i32, align 4
  store i32 %0, i32* %2, align 4
  %3 = load i32, i32* %2, align 4
  call void @exit(i32 %3) #9
  unreachable
}

; Function Attrs: noreturn
declare void @exit(i32) #6

attributes #0 = { noinline nounwind optnone ssp uwtable "darwin-stkchk-strong-link" "disable-tail-calls"="false" "frame-pointer"="all" "less-precise-fpmad"="false" "min-legal-vector-width"="0" "no-infs-fp-math"="false" "no-jump-tables"="false" "no-nans-fp-math"="false" "no-signed-zeros-fp-math"="false" "no-trapping-math"="true" "probe-stack"="___chkstk_darwin" "stack-protector-buffer-size"="8" "target-cpu"="penryn" "target-features"="+cx16,+cx8,+fxsr,+mmx,+sahf,+sse,+sse2,+sse3,+sse4.1,+ssse3,+x87" "tune-cpu"="generic" "unsafe-fp-math"="false" "use-soft-float"="false" }
attributes #1 = { allocsize(0) "darwin-stkchk-strong-link" "disable-tail-calls"="false" "frame-pointer"="all" "less-precise-fpmad"="false" "no-infs-fp-math"="false" "no-nans-fp-math"="false" "no-signed-zeros-fp-math"="false" "no-trapping-math"="true" "probe-stack"="___chkstk_darwin" "stack-protector-buffer-size"="8" "target-cpu"="penryn" "target-features"="+cx16,+cx8,+fxsr,+mmx,+sahf,+sse,+sse2,+sse3,+sse4.1,+ssse3,+x87" "tune-cpu"="generic" "unsafe-fp-math"="false" "use-soft-float"="false" }
attributes #2 = { nounwind "darwin-stkchk-strong-link" "disable-tail-calls"="false" "frame-pointer"="all" "less-precise-fpmad"="false" "no-infs-fp-math"="false" "no-nans-fp-math"="false" "no-signed-zeros-fp-math"="false" "no-trapping-math"="true" "probe-stack"="___chkstk_darwin" "stack-protector-buffer-size"="8" "target-cpu"="penryn" "target-features"="+cx16,+cx8,+fxsr,+mmx,+sahf,+sse,+sse2,+sse3,+sse4.1,+ssse3,+x87" "tune-cpu"="generic" "unsafe-fp-math"="false" "use-soft-float"="false" }
attributes #3 = { nofree nosync nounwind readnone speculatable willreturn }
attributes #4 = { "darwin-stkchk-strong-link" "disable-tail-calls"="false" "frame-pointer"="all" "less-precise-fpmad"="false" "no-infs-fp-math"="false" "no-nans-fp-math"="false" "no-signed-zeros-fp-math"="false" "no-trapping-math"="true" "probe-stack"="___chkstk_darwin" "stack-protector-buffer-size"="8" "target-cpu"="penryn" "target-features"="+cx16,+cx8,+fxsr,+mmx,+sahf,+sse,+sse2,+sse3,+sse4.1,+ssse3,+x87" "tune-cpu"="generic" "unsafe-fp-math"="false" "use-soft-float"="false" }
attributes #5 = { nofree nosync nounwind willreturn }
attributes #6 = { noreturn "darwin-stkchk-strong-link" "disable-tail-calls"="false" "frame-pointer"="all" "less-precise-fpmad"="false" "no-infs-fp-math"="false" "no-nans-fp-math"="false" "no-signed-zeros-fp-math"="false" "no-trapping-math"="true" "probe-stack"="___chkstk_darwin" "stack-protector-buffer-size"="8" "target-cpu"="penryn" "target-features"="+cx16,+cx8,+fxsr,+mmx,+sahf,+sse,+sse2,+sse3,+sse4.1,+ssse3,+x87" "tune-cpu"="generic" "unsafe-fp-math"="false" "use-soft-float"="false" }
attributes #7 = { allocsize(0) }
attributes #8 = { nounwind }
attributes #9 = { noreturn }

!llvm.module.flags = !{!0, !1, !2}
!llvm.ident = !{!3}

!0 = !{i32 2, !"SDK Version", [2 x i32] [i32 12, i32 0]}
!1 = !{i32 1, !"wchar_size", i32 4}
!2 = !{i32 7, !"PIC Level", i32 2}
!3 = !{!"Apple clang version 13.0.0 (clang-1300.0.29.3)"}
!4 = distinct !{!4, !5}
!5 = !{!"llvm.loop.mustprogress"}
