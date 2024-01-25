import time
import sqlite3
from dbPath import findDBPath
import numpy as np
import matplotlib.pyplot as plt
import calcGiniCoefficient as cgc


con = sqlite3.connect(findDBPath("controll.db"))

# lets get gini coefficinet of 16 nodes static network 
# (runID 1 for me) for 2000 rounds

startRound = 350667
endRound = 358667

startRound = 1
endRound = 10000

roundsSavedFor = 100

gini1 = []


gini2 = []
start = time.time()
for i in range(startRound, endRound//roundsSavedFor):
    print(i)
    gini2.append(cgc.caclGiniSql(con, i))
end = time.time()
runTime2 = end-start
print(runTime2)

x = []
for i in range(99):
    x.append((i)*100)
print(len(gini2))
print(gini2)

plt.plot(x,gini2)
plt.xlabel('Round')
plt.ylabel('Gini')
# Set the y-axis ticks with 0 to 1 in increments of 0.05
plt.yticks(np.arange(0.4, 1.05, 0.05))
plt.xticks(np.arange(0, 10500, 500))

# Add horizontal grid lines
grid_values = np.arange(0.4, 1.05, 0.05)
for value in grid_values:
    plt.axhline(y=value, color='gray', linestyle='--', linewidth=0.5)

plt.title('2048 nodes 10000 rounds')


plt.show()




