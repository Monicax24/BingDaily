// Model for Community data matching backend structure.
import 'dart:convert';

class Community {
  final int communityId;
  final String picture;
  final String description;
  final List<int> members;
  final List<int> moderators;
  final List<int> posts;
  final String postTime;
  final String defaultPrompt;

  Community({
    required this.communityId,
    required this.picture,
    required this.description,
    required this.members,
    required this.moderators,
    required this.posts,
    required this.postTime,
    required this.defaultPrompt,
  });

  // Creates a copy of the Community with optional overrides.
  Community copyWith({
    int? communityId,
    String? picture,
    String? description,
    List<int>? members,
    List<int>? moderators,
    List<int>? posts,
    String? postTime,
    String? defaultPrompt,
  }) {
    return Community(
      communityId: communityId ?? this.communityId,
      picture: picture ?? this.picture,
      description: description ?? this.description,
      members: members ?? this.members,
      moderators: moderators ?? this.moderators,
      posts: posts ?? this.posts,
      postTime: postTime ?? this.postTime,
      defaultPrompt: defaultPrompt ?? this.defaultPrompt,
    );
  }

  // Converts Community to a Map for JSON serialization.
  Map<String, dynamic> toMap() {
    return {
      'community_id': communityId,
      'picture': picture,
      'description': description,
      'members': members,
      'moderators': moderators,
      'posts': posts,
      'post_time': postTime,
      'default_prompt': defaultPrompt,
    };
  }

  // Creates Community from a Map (e.g., from JSON).
  factory Community.fromMap(Map<String, dynamic> map) {
    return Community(
      communityId: map['community_id'] as int,
      picture: map['picture'] as String,
      description: map['description'] as String,
      members: List<int>.from(map['members']),
      moderators: List<int>.from(map['moderators']),
      posts: List<int>.from(map['posts']),
      postTime: map['post_time'] as String,
      defaultPrompt: map['default_prompt'] as String,
    );
  }

  // Converts Community to JSON string.
  String toJson() => json.encode(toMap());

  // Creates Community from JSON string.
  factory Community.fromJson(String source) =>
      Community.fromMap(json.decode(source));
}
