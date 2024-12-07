# Introduction to DCDQL(Discount Condition Definition Query Language)

```plaintext
Condition := Expression (LogicalOperator Expression)*
Expression := Key ComparisonOperator Value
Key := "min_cart_value" | "total_price" | "product_id" | "total_category_price" | "category_id"
ComparisonOperator := "==" | "!=" | ">" | "<" | ">=" | "<=" | "in"
Value := Float | String | []String | []Float
LogicalOperator := "AND" | "OR"
Action := "THEN" DiscountType AssignmentOperator Value
DiscountType := "Percentage" | "FlatAmount" | "BOGO"
AssignmentOperator := "="
```

## Syntax Examples

```plaintext
1. min_cart_value > 100 AND product_category == "dairy" THEN Percentage=10
2. product_id IN [1,2,3] OR total_category_price >= 50 THEN FlatAmount=5
```
