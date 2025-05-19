# @@docs Basic Arithmetic
# This section covers basic arithmetic operations.
# The functions here are simple but fundamental.
def add(a, b):
    return a + b

def subtract(a, b):
    return a - b

    # @@docs Advanced Arithmetic
    # This is a nested section with more complex math.
    # Note the indentation - this section is "inside" the previous one.
    def multiply(a, b):
        return a * b

# @@docs Helpers
# Utility functions for the calculator.
def validate_input(x):
    return isinstance(x, (int, float))