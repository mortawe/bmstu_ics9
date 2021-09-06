//  Закаренная функция digits: Int => (Int => List[Int]),
//  выполняющая перевод числа в заданную систему счисления (параметр
//  функции – основание системы счисления).

val digits: Int => (Int => List[Int]) = {
  base =>
    x =>
      if (x < base) List(x) else
        digits(base)(x / base) ::: List(x % base)
}
