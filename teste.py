from typing import Tuple, List
values: List[Tuple[str, str, int]] = []

with open("/tmp/resumo.txt", "r") as file:
    lines = file.readlines()


print(lines)

credit = 1_000 * 100
acumulate_balance = 0
for line in lines :    
    stripedLine =  line.strip()

    if len(stripedLine) == 0:
        continue

    splLine = stripedLine.split(":")

    desc = splLine[0].strip()
    op = splLine[1].strip()
    value = int(splLine[2].strip())
    if op == "d":
        value *= -1

    temp_balance = acumulate_balance + value

    if temp_balance + credit < 0:
        print("Estourou o saldo. ", acumulate_balance, " - ", temp_balance)
        continue

    acumulate_balance += value 

print(acumulate_balance)

    