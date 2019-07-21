-- phpMyAdmin SQL Dump
-- version 4.9.0.1
-- https://www.phpmyadmin.net/
--
-- 主机： 127.0.0.1
-- 生成日期： 2019-07-19 12:48:11
-- 服务器版本： 10.3.16-MariaDB
-- PHP 版本： 7.3.7

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET AUTOCOMMIT = 0;
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- 数据库： `sora_v1`
--

-- --------------------------------------------------------

--
-- 表的结构 `anime`
--

CREATE TABLE `anime` (
  `id` bigint(20) NOT NULL,
  `s_id` bigint(100) NOT NULL DEFAULT 0,
  `name` varchar(255) NOT NULL DEFAULT '""',
  `em_num` varchar(100) NOT NULL,
  `year` varchar(100) NOT NULL DEFAULT '2019',
  `area` varchar(100) NOT NULL DEFAULT '""',
  `picture` varchar(255) NOT NULL DEFAULT '""',
  `introduction` text NOT NULL,
  `form` varchar(100) NOT NULL COMMENT '爬取网站',
  `create_time` timestamp NOT NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- 表的结构 `director`
--

CREATE TABLE `director` (
  `id` int(11) NOT NULL,
  `name` varchar(255) NOT NULL,
  `anime_id` int(100) NOT NULL,
  `create_time` timestamp NOT NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- 表的结构 `drama`
--

CREATE TABLE `drama` (
  `id` int(11) NOT NULL,
  `name` varchar(100) NOT NULL DEFAULT '""',
  `play_url` varchar(255) NOT NULL DEFAULT '""',
  `source` varchar(100) NOT NULL DEFAULT '“”',
  `anime_id` int(100) NOT NULL,
  `create_time` timestamp NOT NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- 表的结构 `star`
--

CREATE TABLE `star` (
  `id` int(11) NOT NULL,
  `name` varchar(255) NOT NULL,
  `anime_id` int(100) NOT NULL,
  `create_time` timestamp NOT NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- 转储表的索引
--

--
-- 表的索引 `anime`
--
ALTER TABLE `anime`
  ADD PRIMARY KEY (`id`);

--
-- 表的索引 `director`
--
ALTER TABLE `director`
  ADD PRIMARY KEY (`id`);

--
-- 表的索引 `drama`
--
ALTER TABLE `drama`
  ADD PRIMARY KEY (`id`);

--
-- 表的索引 `star`
--
ALTER TABLE `star`
  ADD PRIMARY KEY (`id`);

--
-- 在导出的表使用AUTO_INCREMENT
--

--
-- 使用表AUTO_INCREMENT `anime`
--
ALTER TABLE `anime`
  MODIFY `id` bigint(20) NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `director`
--
ALTER TABLE `director`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `drama`
--
ALTER TABLE `drama`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- 使用表AUTO_INCREMENT `star`
--
ALTER TABLE `star`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
