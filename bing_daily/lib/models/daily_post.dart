// Model for Daily Post data matching backend structure.
import 'dart:convert';

class DailyPost {
  final int postId;
  final int communityId;
  final String picture;
  final String caption;
  final int author;
  final DateTime timePosted;
  final int likes;

  DailyPost({
    required this.postId,
    required this.communityId,
    required this.picture,
    required this.caption,
    required this.author,
    required this.timePosted,
    required this.likes,
  });

  // Creates a copy of the DailyPost with optional overrides.
  DailyPost copyWith({
    int? postId,
    int? communityId,
    String? picture,
    String? caption,
    int? author,
    DateTime? timePosted,
    int? likes,
  }) {
    return DailyPost(
      postId: postId ?? this.postId,
      communityId: communityId ?? this.communityId,
      picture: picture ?? this.picture,
      caption: caption ?? this.caption,
      author: author ?? this.author,
      timePosted: timePosted ?? this.timePosted,
      likes: likes ?? this.likes,
    );
  }

  // Converts DailyPost to a Map for JSON serialization.
  Map<String, dynamic> toMap() {
    return {
      'post_id': postId,
      'community_id': communityId,
      'picture': picture,
      'caption': caption,
      'author': author,
      'time_posted': timePosted.toIso8601String(),
      'likes': likes,
    };
  }

  // Creates DailyPost from a Map (e.g., from JSON).
  factory DailyPost.fromMap(Map<String, dynamic> map) {
    return DailyPost(
      postId: map['post_id'] as int,
      communityId: map['community_id'] as int,
      picture: map['picture'] as String,
      caption: map['caption'] as String,
      author: map['author'] as int,
      timePosted: DateTime.parse(map['time_posted']),
      likes: map['likes'] as int,
    );
  }

  // Converts DailyPost to JSON string.
  String toJson() => json.encode(toMap());

  // Creates DailyPost from JSON string.
  factory DailyPost.fromJson(String source) =>
      DailyPost.fromMap(json.decode(source));
}
