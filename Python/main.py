import os
def main():
    cities: dict[str, tuple[float, int]] = {}
    with os.open("../measurements.txt") as file:
        file: list[str]
        for line in file:
            city, temp = line.strip().split(';')
            count = cities[city][1]
            tempAvg = cities[city][0] * count
            count += 1
            cities[city][0] = (tempAvg + temp) / count
            cities[city][1] = count
main()