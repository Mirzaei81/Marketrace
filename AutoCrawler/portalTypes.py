from enum import Enum
from dataclasses import dataclass
from typing import Optional, List, Any, TypeVar, Callable, Type, cast


T = TypeVar("T")
EnumT = TypeVar("EnumT", bound=Enum)


def from_int(x: Any) -> int:
    assert isinstance(x, int) and not isinstance(x, bool)
    return x


def from_str(x: Any) -> str:
    assert isinstance(x, str)
    return x
def from_strOrNone(x:Any)->str:
    if x==None:
        return ""
    else:
        assert isinstance(x, str)
        return x

def from_none(x: Any) -> Any:
    assert x is None
    return x


def from_union(fs, x):
    for f in fs:
        try:
            return f(x)
        except:
            pass
    assert False


def from_list(f: Callable[[Any], T], x: Any) -> List[T]:
    assert isinstance(x, list)
    return [f(y) for y in x]


def from_float(x: Any) -> float:
    assert isinstance(x, (float, int)) and not isinstance(x, bool)
    return float(x)


def is_type(t: Type[T], x: Any) -> T:
    assert isinstance(x, t)
    return x


def to_enum(c: Type[EnumT], x: Any) -> EnumT:
    assert isinstance(x, c)
    return x.value


def to_float(x: Any) -> float:
    assert isinstance(x, (int, float))
    return x


def from_bool(x: Any) -> bool:
    assert isinstance(x, bool)
    return x


def to_class(c: Type[T], x: Any) -> dict:
    assert isinstance(x, c)
    return cast(Any, x).to_dict()


class Status(Enum):
    APPROVED = "approved"
    BANK_PAYMENT = "bank_payment"
    CASH_ON_DELIVERY = "cash_on_delivery"
    ONLINE_PAYMENT = "online_payment"
    SHIPPING_REQUIRED = "shipping_required"


class TypeEnum(Enum):
    COMMODITY = "commodity"


@dataclass
class Variant:
    id: int
    product_id: int
    title: str
    tax: None
    shipping: None
    length: None
    width: None
    height: None
    sku: Optional[int] = None
    image: None = None
    type: TypeEnum = TypeEnum.COMMODITY
    status: List[Status]  = None
    files: None =None
    price: Optional[int] = None
    compare_price: Optional[int] = None
    weight: Optional[int] = None
    stock: Optional[float] = None
    minimum: Optional[int] = None
    maximum: Optional[int] = None

    @staticmethod
    def from_dict(obj: Any) -> 'Variant':
        assert isinstance(obj, dict)
        id = from_int(obj.get("id"))
        product_id = from_int(obj.get("product_id"))
        title = from_str(obj.get("title"))
        tax = from_none(obj.get("tax"))
        shipping = from_none(obj.get("shipping"))
        length = from_none(obj.get("length"))
        width = from_none(obj.get("width"))
        height = from_none(obj.get("height"))
        sku = from_union([from_none, lambda x: int(from_str(x))], obj.get("sku"))
        image = from_strOrNone(obj.get("image"))
        type = TypeEnum(obj.get("type"))
        status = from_list(Status, obj.get("status"))
        files = from_none(obj.get("files"))
        price = from_union([from_none, from_int], obj.get("price"))
        compare_price = from_union([from_none, from_int], obj.get("compare_price"))
        weight = from_union([from_none, from_int], obj.get("weight"))
        stock = from_union([from_none, from_float], obj.get("stock"))
        minimum = from_union([from_none, from_int], obj.get("minimum"))
        maximum = from_union([from_none, from_int], obj.get("maximum"))
        return Variant(id, product_id, title, tax, shipping, length, width, height, sku, image, type, status, files, price, compare_price, weight, stock, minimum, maximum)

    def to_dict(self) -> dict:
        result: dict = {}
        result["id"] = from_int(self.id)
        result["product_id"] = from_int(self.product_id)
        result["title"] = from_str(self.title)
        result["tax"] = from_none(self.tax)
        result["shipping"] = from_none(self.shipping)
        result["length"] = from_none(self.length)
        result["width"] = from_none(self.width)
        result["height"] = from_none(self.height)
        result["sku"] = from_union([lambda x: from_none((lambda x: is_type(type(None), x))(x)), lambda x: from_str((lambda x: str((lambda x: is_type(int, x))(x)))(x))], self.sku)
        result["image"] = from_none(self.image)
        result["type"] = to_enum(TypeEnum, self.type)
        result["status"] = from_list(lambda x: to_enum(Status, x), self.status)
        result["files"] = from_none(self.files)
        result["price"] = from_union([from_none, from_int], self.price)
        result["compare_price"] = from_union([from_none, from_int], self.compare_price)
        result["weight"] = from_union([from_none, from_int], self.weight)
        result["stock"] = from_union([from_none, to_float], self.stock)
        result["minimum"] = from_union([from_none, from_int], self.minimum)
        result["maximum"] = from_union([from_none, from_int], self.maximum)
        return result


@dataclass
class VariantResult:
    success: bool
    total: int
    count: int
    variants: List[Variant]

    @staticmethod
    def from_dict(obj: Any) -> 'VariantResult':
        assert isinstance(obj, dict)
        success = from_bool(obj.get("success"))
        total = from_int(obj.get("total"))
        count = from_int(obj.get("count"))
        variants = from_list(Variant.from_dict, obj.get("variants"))
        return VariantResult(success, total, count, variants)

    def to_dict(self) -> dict:
        result: dict = {}
        result["success"] = from_bool(self.success)
        result["total"] = from_int(self.total)
        result["count"] = from_int(self.count)
        result["variants"] = from_list(lambda x: to_class(Variant, x), self.variants)
        return result


def variant_from_dict(s: Any) -> VariantResult:
    return VariantResult.from_dict(s)


def variant_to_dict(x: VariantResult) -> Any:
    return to_class(VariantResult, x)
