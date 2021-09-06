/* Класс Permutation[T], представляющий неизменяемую перестановку
подмножества значений типа T с двумя операциями: получение i-го
элемента и транспозиция. В случае, если T – строковый тип, для
Permutation[T] доступна дополнительная операция superStringSize,
возвращающая длину минимальной суперстроки, в которую строки,
содержащиеся в перестановке, входят в том порядке, в котором они
расположены в перестановке.*/

abstract class MinSuperFinder[T] {
  def findMinSuper(list: List[T]): T

  def findMinSuperSize(list: List[T]): Int
}

object MinSuperFinder {

  implicit object str extends MinSuperFinder[String] {
    def maxEndsWith(len: Int, str: String, prefix: String, maxLen: Int): Int = {
      if (prefix.length >= len) {
        if (str.endsWith(prefix.substring(0, len))) {
          maxEndsWith(len + 1, str, prefix, len)
        } else {
          maxEndsWith(len + 1, str, prefix, maxLen)
        }
      } else maxLen
    }

    def superStr(prefix: String, tail: List[String]): String = {
      if (tail == Nil) prefix else {
        val currentStr :: newTail = tail
        val prefixIntersection = maxEndsWith(1, prefix, currentStr, 0)
        superStr(prefix + currentStr.substring(prefixIntersection), newTail)
      }
    }

    override def findMinSuper(list: List[String]): String = {
      superStr("", list)
    }

    override def findMinSuperSize(list: List[String]): Int = {
      findMinSuper(list).length
    }

  }

}

class Permutation[T](elems: List[T]) {
  val e = elems

  def getByPos(i: Int): T = {
    e(i)
  }

  def transpose(i: Int, j: Int): Permutation[T] = {
    val ith = e(i)
    val jth = e(j)
    val updatedI = e.updated(i, jth)
    new Permutation[T](updatedI.updated(j, ith))
  }

  def SuperStringSize()(implicit minSuperStrFinder: MinSuperFinder[T]): Int = {
    minSuperStrFinder.findMinSuperSize(e)
  }
}

val numPerm = new Permutation[Int](List(1, 2, 3, 4, 5))
println(numPerm.getByPos(1))
println(numPerm.transpose(1, 3).e)
//println(numPerm.superStringSize()) not works for non-String


val stringPerm = new Permutation[String](List("abba", "bba", "ba", "a"))
println(stringPerm.SuperStringSize())
