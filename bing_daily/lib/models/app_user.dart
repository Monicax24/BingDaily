import 'dart:convert';

import 'package:flutter/foundation.dart';

// ignore_for_file: public_member_api_docs, sort_constructors_first
class AppUser {
  final String userId;
  final String email;
  final String username;
  final String joinDate;
  List<dynamic> communities;

  AppUser({
    required this.userId,
    required this.email,
    required this.username,
    required this.joinDate,
    required this.communities,
  });

  AppUser copyWith({
    String? userId,
    String? email,
    String? username,
    String? joinDate,
    List<dynamic>? communities,
  }) {
    return AppUser(
      userId: userId ?? this.userId,
      email: email ?? this.email,
      username: username ?? this.username,
      joinDate: joinDate ?? this.joinDate,
      communities: communities ?? this.communities,
    );
  }

  Map<String, dynamic> toMap() {
    return <String, dynamic>{
      'userId': userId,
      'email': email,
      'username': username,
      'joinDate': joinDate,
      'communities': communities,
    };
  }

  factory AppUser.fromMap(Map<String, dynamic> map) {
    return AppUser(
      userId: map['userId'] as String,
      email: map['email'] as String,
      username: map['username'] as String,
      joinDate: map['joinDate'] as String,
      communities: map['communities'],
    );
  }

  String toJson() => json.encode(toMap());

  factory AppUser.fromJson(String source) =>
      AppUser.fromMap(json.decode(source) as Map<String, dynamic>);

  @override
  String toString() {
    return 'AppUser(userId: $userId, email: $email, username: $username, joinDate: $joinDate, communities: $communities)';
  }

  @override
  bool operator ==(covariant AppUser other) {
    if (identical(this, other)) return true;

    return other.userId == userId &&
        other.email == email &&
        other.username == username &&
        other.joinDate == joinDate &&
        listEquals(other.communities, communities);
  }

  @override
  int get hashCode {
    return userId.hashCode ^
        email.hashCode ^
        username.hashCode ^
        joinDate.hashCode ^
        communities.hashCode;
  }
}
