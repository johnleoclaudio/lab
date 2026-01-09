export function march(
  user: User,
  piece: Piece,
  publishMove: (move: Move) => void,
): void {
  // do this
}

export type User = {
  name: string;
  pieces: Piece[];
};

export type Move = {
  username: string;
  piece: Piece;
};

export type Piece = {
  location: string;
  name: string;
};

export function doBattles(publishMoves: Move[], users: User[]): Piece[] {
  const fights: Piece[] = [];
  for (const move of publishMoves) {
    for (const user of users) {
      if (move.username === user.name) {
        continue;
      }
      for (const piece of user.pieces) {
        if (piece.location == move.piece.location) {
          fights.push(piece);
        }
      }
    }
  }
  return fights;
}

