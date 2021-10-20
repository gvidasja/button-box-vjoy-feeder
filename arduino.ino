byte encoderPins[] = {A0, A1, A2, A3, A4, A5, A6, A7};
const int encCount = sizeof(encoderPins) / sizeof(encoderPins[0]);
byte encoderSent[encCount];              // offset bottom encoders, so that input is trigger on the 'click' (Y, G, R, B)
const byte encoderSkip[] = {2, 2, 1, 1}; // skips nth signal from encoder (Y, G, R, B)

byte rows[] = {2, 3, 4, 5, 6};
const int rowCount = sizeof(rows) / sizeof(rows[0]);

byte cols[] = {7, 8, 9, 10, 11};
const int colCount = sizeof(cols) / sizeof(cols[0]);

byte encoders[8];
byte prevEncoders[8];
byte keys[colCount][rowCount];
byte prevKeys[colCount][rowCount];
byte debounce[colCount][rowCount];

byte DEBOUNCE_DURATION = 200;

void send(bool msg, int button)
{
  Serial.print(msg);
  Serial.println(button);
}

void setup()
{
  Serial.begin(9600);

  for (int x = 0; x < rowCount; x++)
  {
    pinMode(rows[x], INPUT);
  }

  for (int x = 0; x < colCount; x++)
  {
    pinMode(cols[x], INPUT_PULLUP);
  }

  for (int i = 0; i < encCount; i += 2)
  {
    pinMode(encoderPins[i], INPUT);
    pinMode(encoderPins[i + 1], INPUT);
    encoderSent[i / 2] = 0;
  }

  auto time = millis();

  for (int row = 0; row < rowCount; row++)
  {
    for (int col; col < colCount; col++)
    {
      debounce[col][row] = time;
    }
  }
}

void readMatrix()
{
  for (int colIndex = 0; colIndex < colCount; colIndex++)
  {
    byte curCol = cols[colIndex];
    pinMode(curCol, OUTPUT);
    digitalWrite(curCol, LOW);

    for (int rowIndex = 0; rowIndex < rowCount; rowIndex++)
    {
      byte rowCol = rows[rowIndex];
      pinMode(rowCol, INPUT_PULLUP);
      prevKeys[colIndex][rowIndex] = keys[colIndex][rowIndex];
      keys[colIndex][rowIndex] = digitalRead(rowCol);
      pinMode(rowCol, INPUT);
    }

    pinMode(curCol, INPUT);
  }

  for (int i = 0; i < encCount; i++)
  {
    byte cur = encoderPins[i];
    pinMode(cur, INPUT_PULLUP);
    prevEncoders[i] = encoders[i];
    encoders[i] = digitalRead(cur);
  }
}

void sendMatrix()
{
  for (int rowIndex = 0; rowIndex < rowCount; rowIndex++)
  {
    for (int colIndex = 0; colIndex < colCount; colIndex++)
    {
      if (keys[colIndex][rowIndex] < prevKeys[colIndex][rowIndex])
      {
        auto time = millis();

        if (time - debounce[colIndex][rowIndex] > DEBOUNCE_DURATION)
        {
          send(true, colIndex * rowCount + rowIndex);
          debounce[colIndex][rowIndex] = time;
        }
      }

      if (keys[colIndex][rowIndex] > prevKeys[colIndex][rowIndex])
      {
        send(false, colIndex * rowCount + rowIndex);
      }
    }
  }

  for (int i = 0; i < encCount; i += 2)
  {
    if ((encoders[i] != prevEncoders[i]))
    {
      if (encoderSent[i / 2] == 0)
      {
        send(true, (colCount * rowCount) + i + (encoders[i] == encoders[i + 1]));
        send(false, (colCount * rowCount) + i + (encoders[i] == encoders[i + 1]));
      }
      encoderSent[i / 2] = (encoderSent[i / 2] + 1) % encoderSkip[i / 2];
    }
  }
}

long long a = 0;

void loop()
{
  readMatrix();
  sendMatrix();
}
