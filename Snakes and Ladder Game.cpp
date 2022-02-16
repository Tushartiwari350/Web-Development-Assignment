#include<iostream>
#include<stdlib.h>
#include<stdio.h>
#include<time.h>
#include<conio.h>
#include<vector>
#include<queue>
#include<algorithm>

using namespace std;

//Global Board Variables
int n = 10, m = 10;
vector<int> jump(n * m  + 7, 0);
vector<int> adjl[101];
int num_of_ladders = 8, num_of_snakes = 13;

void draw_line(int n, char ch);                                 //printing horizontal lines for presentation
void board();                                                   //printing location of snakes and ladders on the board
void gamescore(char name1[], char name2[], int p1, int p2);     //printing player locations and shortes past
void play_dice(int &score);                                     //random no. from 1-6
void set_jump();
void set_adjl();

int main()  {
    int player1 = 1, player2 = 1, lastposition;
    char player1name[80], player2name[80];
    set_jump();
    set_adjl();

    system("cls");
    srand(time(NULL));
    draw_line(50, '=');
    cout << "\n\n\n\n\t\tSNAKE LADDER GAME\n\n\n\n";
    draw_line(50, '=');

    cout << "\n\n\nEnter name of player 1: ";
    cin >> player1name;
    cout << "\nEnter name of player 2: ";
    cin >> player2name;

    getch();
    board();                                          //showing players the board before the game starts

    while (player1 <= 100 && player2 <= 100)
    {   system("cls");
        gamescore(player1name, player2name, player1, player2);
        cout << "\n\n--->" << player1name << "'s turn, \tPress any key to roll";
        getch();
        lastposition = player1;
        play_dice(player1);

        if (player1 < lastposition)                     //if player's movement is backwards => snake encountered
            cout << "\nYOU LANDED ON A SNAKE!! You are now at postion " << player1 << "\n";

        else if (player1 > lastposition + 6)            //if player's movement is forward(>6) => ladder encountered
            cout << "\nYOU LANDED ON A LADDER!! You are now at position " << player1;


        cout << "\n\n--->" << player2name << "'s turn, \tPress any key to roll ";

        getch();
        lastposition = player2;
        play_dice(player2);

        if (player2 < lastposition)
            cout << "\nYOU LANDED ON A SNAKE!! You are now at postion " << player2 << "\n";

        else if (player2 > lastposition + 6)
            cout << "\nYOU LANDED ON A LADDER!! You are now at position " << player2 << "\n";
        getch();
    }

    system("cls");
    //GAME END PAGE
    cout << "\n\n\n";
    draw_line(50, '*');
    cout << "\n\n\t\t      RESULT\n\n";
    draw_line(50, '*');
    cout << endl;
    gamescore(player1name, player2name, player1, player2);
    cout << "\n\n\n";

    if (player1 >= player2)
        cout << player1name << " WINS!\n\n";
    else
        cout << player2name << " WINS!\n\n";
    draw_line(50, '-');
    getch();

}


void draw_line(int n, char ch) {
    for (int i = 0; i < n; i++) cout << ch;
}

void board()  {
    system("cls");
    cout << "\n\n";
    draw_line(50, '-');
    cout << "\n\t\tSNAKE POSTIONS\n";
    draw_line(50, '-');
    cout << "\n\tFrom 98 to 28 \n\tFrom 95 to 24\n\tFrom 92 to 51\n\tFrom 83 to 19\n\tFrom 73 to  1\n\tFrom 69 to 33\n\tFrom 64 to 36\n\tFrom 59 to 17\n\tFrom 55 to  7\n\tFrom 52 to 11\n\tFrom 48 to  9\n\tFrom 46 to  5\n\tFrom 44 to 22\n\n";
    draw_line(50, '-');
    cout << "\n\t\t LADDER POSITIONS\n";
    draw_line(50, '-');
    cout << "\n\tFrom  8 to 26\n\tFrom 21 to 82\n\tFrom 43 to 77\n\tFrom 50 to 91\n\tFrom 62 to 96\n\tFrom 66 to 87\n\tFrom 80 to 100\n";
    draw_line(50, '-');
    cout << endl;

}

void play_dice(int &score)  {
    int dice;
    dice = (rand() % 6) + 1;                          //simulating a dice => 1 to 6
    cout << "\nYou rolled a " << dice;
    score = score + dice;
    cout << endl << "Now you are at position " << score;
    switch (score)                                    //special cases(i.e. snake/ladder)
    {
    case 98 : score = 28; break;
    case 95 : score = 24; break;
    case 92 : score = 51; break;
    case 83 : score = 19; break;
    case 73 : score = 1; break;
    case 69 : score = 33; break;
    case 64 : score = 36; break;
    case 59 : score = 17; break;
    case 55 : score = 7; break;
    case 52 : score = 11; break;
    case 48 : score = 9; break;
    case 46 : score = 5; break;
    case 44 : score = 22; break;
    case 8  : score = 26; break;
    case 21 : score = 82; break;
    case 43 : score = 77; break;
    case 50 : score = 91; break;
    case 54 : score = 93; break;
    case 62 : score = 96; break;
    case 66 : score = 87; break;
    case 80 : score = 100;
    }
    cin.get();
}

void set_jump() {
    //ladder jumps
    jump[8] = 26 - 8;
    jump[21] = 82 - 21;
    jump[43] = 77 - 43;
    jump[50] = 91 - 50;
    jump[54] = 93 - 54;
    jump[62] = 96 - 62;
    jump[66] = 87 - 66;
    jump[80] = 100 - 80;

    //snake jumps
    jump[98] = 28 - 98;
    jump[95] = 24 - 95;
    jump[92] = 51 - 92;
    jump[83] = 19 - 83;
    jump[73] = 1 - 73;
    jump[69] = 33 - 69;
    jump[64] = 36 - 64;
    jump[59] = 17 - 59;
    jump[55] = 7 - 55;
    jump[52] = 11 - 52;
    jump[48] = 9 - 48;
    jump[46] = 5 - 46;
    jump[44] = 22 - 44;
}

void set_adjl() {
    for (int i = 1; i <= n * m; i++) {
        for (int j = 1; j < 7; j++) {
            int cur = i + j;
            cur += jump[cur];
            if (cur <= n * m) {
                adjl[i].push_back(cur);
            }
        }
    }
}

void solve(int pos) {
    //solve function starts
    vector<int> par(n * m + 7);
    vector<int> visited(n * m + 7, 0);

    visited[pos] = 1;

    queue<int> q;

    q.push(pos);

    while (!q.empty()) {
        int cur_block = q.front();
        q.pop();
        for (auto x : adjl[cur_block]) {
            if (!visited[x]) {
                visited[x] = visited[cur_block] + 1;
                par[x] = cur_block;
                q.push(x);
            }
        }
        if (visited[n * m]) {
            break;
        }
    }

    cout << "Dice throws required= " << visited[n * m] - 1;
    cout << endl;
    cout << "Shortest path -> ";
    vector<int> ans;
    int num = n * m;
    while (num != pos) {
        ans.push_back(num);
        num = par[num];
    }
    ans.push_back(pos);
    reverse(ans.begin(), ans.end());
    for (int i = 0; i < ans.size(); i++) {
        cout << ans[i] << " ";
    }
    cout << endl;
}

void gamescore(char name1[], char name2[], int p1, int p2)  {
    cout << "\n";
    draw_line(50, '.');
    cout << "\n\t\t   GAME STATUS\n";
    draw_line(50, '.');
    cout << endl << "\t--->" << name1 << "'s position: " << p1 << endl;
    cout << "\t--->" << name2 << "'s position: " << p2 << endl;
    if ((p1 < 100) && (p2 < 100)) {
        cout << "\nShortest path for " << name1 << ":" << endl;
        solve(p1);
        cout << "\nShortest path for " << name2 << ":" << endl;
        solve(p2);
    }

    draw_line(50, '_');
    cout << endl;
}