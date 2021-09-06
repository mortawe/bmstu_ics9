/*Система неравенств вида a1x1 + a2x2 + . . . + aN xN ≤ b. Конструктор
класса должен принимать список коэффициентов ai и свободный член b и
порождать систему, состоящую из одного неравенства. Операции: «+» –
объединение двух систем в одну; «/» – принимает число i и возвращает
систему, полученную из данной системы путём присвоения нулевого
значения i-той переменной; «check» – проверка, удовлетворяет ли список
значений переменных системе неравенств.*/

class Inequality(ax: Vector[Double], bx: Double) {
  var a: Vector[Double] = ax
  var b: Double = bx

  def defaults() {
    a = ax
    b = bx
  }

  def check(x: Vector[Double]): Boolean = {
    val xa = x zip a map { case (x, a) => x * a }
    xa.sum <= b
  }

  def +(o: Inequality): Inequality = {
    val resultA = o.a zip a map { case (a, b) => a + b }
    val resultB = b + o.b
    new Inequality(resultA, resultB)
  }

  def /(i: Int): Inequality = {
    if (i >= a.size) {
      this
    } else {
      val (before, tail) = a.splitAt(i)
      val (_, after) = tail.splitAt(1)
      new Inequality((before :+ 0.0) ++ after, b)
    }
  }
}








